using HiveServer;
using MySql.Data.MySqlClient;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using System.Web;
using System.Web.Mvc;
using XMLParser.Models;
using XMLParser.AuthenticationDTO;
using XMLParser.Filter;
using System.Threading;
using System.Web.Script.Serialization;
using Models.HiveModels;
using XMLParser.Controllers.XMLParser;
using Newtonsoft.Json.Linq;
using Newtonsoft.Json;
using System.Configuration;

namespace XMLParser.Controllers
{
    public class HomeController : Controller
    {
        static string ConnectionString = ConfigurationManager.AppSettings["connectionString"];
        public ActionResult Index()
        {
            if (System.Web.HttpContext.Current.Session["Email"] != null)
            {
                return RedirectToAction("CreateAssetFormView", "Home");
            }
            else
            {
                ViewBag.Title = "Home Page";
                ViewBag.Current = "Home";
                return View();
            }
        }

        public ActionResult HiveView()
        {
            return View();
        }

        [CustomActionFilter]
        [OutputCache(NoStore = true, Duration = 0, VaryByParam = "*")]
        public ActionResult CreateProjectView()
        {
            ViewBag.Current = "DropDown";
            return View();
        }

        [CustomActionFilter]
        [OutputCache(NoStore = true, Duration = 0, VaryByParam = "*")]
        public ActionResult ProjectConfiguration()
        {
            ViewBag.Current = "Project Configuration";
            return View();
        }
        



        [CustomActionFilter]
        [OutputCache(NoStore = true, Duration = 0, VaryByParam = "*")]
        public ActionResult ProjectsView()
        {
            ViewBag.Current = "View Projects";
            return View();
        }

        [CustomActionFilter]
        [OutputCache(NoStore = true, Duration = 0, VaryByParam = "*")]
        public ActionResult AssetsView()
        {
            ViewBag.Current = "View Assets";
            return View();
        }

        [CustomActionFilter]
        [OutputCache(NoStore = true, Duration = 0, VaryByParam = "*")]
        public ActionResult UsersView()
        {
            ViewBag.Current = "View Users";
            return View();
        }

        [CustomActionFilter]
        [OutputCache(NoStore = true, Duration = 0, VaryByParam = "*")]
        public ActionResult TasksView()
        {
            ViewBag.Current = "View Tasks";
            return View();
        }

        [CustomActionFilter]
        [OutputCache(NoStore = true, Duration = 0, VaryByParam = "*")]
        public ActionResult CreateTaskView()
        {
            ViewBag.Current = "DropDown";
            return View();
        }

        [CustomActionFilter]
        [OutputCache(NoStore = true, Duration = 0, VaryByParam = "*")]
        public ActionResult CreateAssignmentView()
        {
            ViewBag.Current = "DropDown";
            return View();
        }

        [CustomActionFilter]
        [OutputCache(NoStore = true, Duration = 0, VaryByParam = "*")]
        public ActionResult MyXMLParser()
        {
            ViewBag.Current = "DropDown";
            return View();
        }

        public ActionResult Signup()
        {
            ViewBag.Current = "Signup";
            return View();
        }

        
        public ActionResult Login()
        {
            ViewBag.Current = "Login";
            System.Web.HttpContext.Current.Session["UserName"] = 10;
            return View();
        }

        public ActionResult AssetCreationModeSelection()
        {
            ViewBag.Current = "AssetCreationModeSelection";
            return View();
        }
        public ActionResult LogOut()
        {
            ViewBag.Current = "Home";
            Session.Abandon();
            return RedirectToAction("Index", "Home");
        }


        //======================================================
        //===========                              =============
        //===========    Server storage testing    =============
        //===========                              =============
        //======================================================
        [CustomActionFilter]
        [OutputCache(NoStore = true, Duration = 0, VaryByParam = "*")]
        public ActionResult serverStorageTesting()
        {
            ViewBag.Current = "serverStorageTesting";
            return View();
        }

        [HttpGet]
        public JsonResult createProject(string number)
        {
            int count = 1;
            int value = Convert.ToInt32(number);

            HiveController hc = new HiveController();
            hc.ControllerContext = new ControllerContext(this.Request.RequestContext, hc);
            while (count <= value)
            {
                string projectId = DateTime.Now.ToOADate() + "Id" + count;
                string projectName = DateTime.Now.ToOADate() + "Name" + count;
                string projectDes = DateTime.Now.ToOADate() + "Des" + count;
                var task = Task.Run(async () => { await hc.CreateProject(projectId, projectName, projectDes); });
                task.Wait();
                count++;
                task.Dispose();
            }

            return Json(new
            {
                aaData = "{ value:ok }"
            }, JsonRequestBehavior.AllowGet);
        }

        [HttpGet]
        public JsonResult createTask(string number, string projectId, string state, int total, int matching)
        {
            int count = 1;
            int value = Convert.ToInt32(number);

            HiveController hc = new HiveController();
            while (count <= value)
            {
                string tProjectId = projectId;
                string tName = DateTime.Now.ToOADate() + "Name" + count;
                string tDes = DateTime.Now.ToOADate() + "Des" + count;
                string tState = state;
                int tTotal = total;
                int tMatching = matching;
                var task = Task.Run(async () => { await hc.CreateTask(tProjectId, tName, tDes, tState, tTotal, tMatching); });
                task.Wait();
                count++;
                task.Dispose();
            }

            return Json(new
            {
                aaData = "{ value:ok }"
            }, JsonRequestBehavior.AllowGet);
        }

        [HttpGet]
        public JsonResult createAssignment(string number, string projectId, string task_id, string asset_ids, string user_id)
        {
            int count = 0;
            int value = Convert.ToInt32(number);
            string[] assetIds = asset_ids.Split('#');

            HiveController hc = new HiveController();
            foreach (var a in assetIds)
            {
                if (a != "")
                {
                    string tProjectId = projectId;
                    string task_Id = task_id;
                    string userId = user_id;

                    var task = Task.Run(async () => { await hc.CreateAssignment(tProjectId, task_Id, a, userId); });
                    task.Wait();
                    count++;
                    task.Dispose();
                }
                if (count == value)
                    break;
            }

            return Json(new
            {
                aaData = "{ value:ok }"
            }, JsonRequestBehavior.AllowGet);
        }

        [HttpGet]
        public JsonResult createAsset(string number, string data)
        {
            JavaScriptSerializer jss = new JavaScriptSerializer();
            ProjectQuestionModel asset = jss.Deserialize<ProjectQuestionModel>(data);
            int count = 0;
            int value = Convert.ToInt32(number);
            JObject json = JObject.Parse(data);
            XMLParserController xpc = new XMLParserController();
            dynamic v = JsonConvert.DeserializeObject(data);
            while (count <= value)
            {
                xpc.GenerateHiveCall(asset);
                count++;

            }


            return Json(new
            {
                aaData = "{ value:ok }"
            }, JsonRequestBehavior.AllowGet);
        }

        //======================================================


        /// <summary>
        /// ActionResult for ordinary session(HttpContext).
        /// </summary>
        /// <param name="sessionValue"></param>
        /// <returns></returns>
        [HttpPost]
        [Route("SaveSession")]
        public ActionResult SaveSession(string sessionValue)
        {
            try
            {
                System.Web.HttpContext.Current.Session["sessionString"] = sessionValue;
                ViewData["sessionString"] = sessionValue;
                return View();
            }
            catch (InvalidOperationException)
            {
                return View();
            }
        }

        [HttpPost]
        [Route("LoginCall")]
        public JsonResult LoginCall(string email,string password)
        {
            bool validUser = false;
            MySqlConnection connection = new MySqlConnection(ConnectionString);
            try
            {
                connection.Open();
                string query = "SELECT * FROM users WHERE (EmailId = '" + email + "' OR UserName = '"+ email + "') AND Password = '" + password + "'";
                MySqlCommand cmd = new MySqlCommand(query, connection);
                using (MySqlDataReader dataReader = cmd.ExecuteReader())
                {
                    if (dataReader.Read())
                    {
                        System.Web.HttpContext.Current.Session["Email"] = dataReader["EmailId"];
                        System.Web.HttpContext.Current.Session["FirstName"] = dataReader["FirstName"];
                        System.Web.HttpContext.Current.Session["UserName"] = dataReader["UserName"];
                        validUser = true;
                    }
                    dataReader.Close();
                }

                if (validUser)
                {
                    string updateLastLoginTime = "update users set LastLoginTime = '" + DateTime.Now.ToString("yyyy-MM-dd HH:mm:ss.fff") + "' WHERE EmailId = '" + email + "' AND Password = '" + password + "'";
                    MySqlCommand updateCmd = new MySqlCommand(updateLastLoginTime, connection);
                    updateCmd.ExecuteNonQuery();
                }
            }
            catch (Exception e)
            {
                throw e;
            }
            finally
            {
                if (connection.State == System.Data.ConnectionState.Open)
                {
                    connection.Close();
                }
            }
            return Json(validUser);
        }

        [HttpGet]
        // GET: UpdateUser/AssociateQuestion
        public ActionResult AssociateQuestion()
        {
            return PartialView("AssociateQuestion");
        }


        [HttpPost]
        [Route("RegisterCall")]
        public JsonResult RegisterCall(RegisterDTO registerdto)
        {
            bool validUser = false;
            MySqlConnection connection = new MySqlConnection(ConnectionString);
            //string DOB = registerdto.DOB.ToString("yyyy-MM-dd");
            try
            {
                string query = "";
                connection.Open();
                if (registerdto != null)
                {
                    query = "SELECT * FROM users WHERE EmailId = '" + registerdto.EmailId + "'";
                    MySqlCommand cmdForcheckingExistence = new MySqlCommand(query, connection);
                    using (MySqlDataReader dataReader = cmdForcheckingExistence.ExecuteReader())
                    {
                        if (dataReader.Read())
                            return Json("User with same email already exists!");

                        dataReader.Close();
                    }
                    registerdto.registrationTime = DateTime.Now;
                    query = "INSERT INTO users(FirstName, LastName, StudentId, EmailId, Password,UserName,RegistrationTime) VALUES ('" +
                         registerdto.FirstName + "','" +
                         registerdto.LastName + "','" +
                         registerdto.StudentId + "','" +
                         registerdto.EmailId + "','" +
                         registerdto.Password + "','" +
                         registerdto.username +"','" +
                         registerdto.registrationTime.ToString("yyyy-MM-dd HH:mm:ss.fff") + "')";
                    MySqlCommand cmdForInsertion = new MySqlCommand(query, connection);
                    int rows = cmdForInsertion.ExecuteNonQuery();
                    if (rows > 0)
                    {
                        validUser = true;
                    }
                }
            }
            catch (Exception)
            {
                throw;
            }
            finally
            {
                if (connection.State == System.Data.ConnectionState.Open)
                {
                    connection.Close();
                }
            }
            return Json(validUser);
        }
        [CustomActionFilter]
        public ActionResult UploadFileView()
        {
            ViewBag.Current = "Upload File View";
            return View();
        }

        [CustomActionFilter]
        [OutputCache(NoStore = true, Duration = 0, VaryByParam = "*")]
        public ActionResult EditProfile()
        {
            MySqlConnection connection = new MySqlConnection(ConnectionString);
            try
            {
                string query = "";
                connection.Open();
                query = "SELECT * FROM users WHERE EmailId = '" + Session["Email"] + "'";
                MySqlCommand cmdForcheckingExistence = new MySqlCommand(query, connection);
                using (MySqlDataReader dataReader = cmdForcheckingExistence.ExecuteReader())
                {
                    if (dataReader.Read())
                    {
                        RegisterDTO user = new RegisterDTO();
                        user.EmailId = dataReader["EmailId"].ToString();
                        user.FirstName = dataReader["FirstName"].ToString();
                        user.LastName = dataReader["LastName"].ToString();
                        user.Password = dataReader["Password"].ToString();
                        user.StudentId = dataReader["StudentId"].ToString();
                        user.username = dataReader["UserName"].ToString();

                        ViewBag.User = user;
                    }
                    dataReader.Close();
                }
            }
            catch (Exception)
            {
                throw;
            }
            finally
            {
                if (connection.State == System.Data.ConnectionState.Open)
                {
                    connection.Close();
                }
            }
            ViewBag.Current = "EditProfile";
            return View();
        }


        [HttpPost]
        [CustomActionFilter]
        [OutputCache(NoStore = true, Duration = 0, VaryByParam = "*")]
        public JsonResult UpdateProfile(RegisterDTO registerdto)
        {
            bool isUpdated = false;
            MySqlConnection connection = new MySqlConnection(ConnectionString);
            //string DOB = registerdto.DOB.ToString("yyyy-MM-dd");
            try
            {
                connection.Open();
                if (registerdto != null)
                {
                    string updateLastLoginTime = "update users set " +
                        "FirstName = '" + registerdto.FirstName + "', " +
                        "LastName ='" + registerdto.LastName + "', " +
                        "Password='" + registerdto.Password + "', " +
                        "StudentId='" + registerdto.StudentId + "', " +
                        "UserName='"+registerdto.username+ "'"+
                        "WHERE EmailId = '" + Session["Email"]+"'";
                    MySqlCommand updateCmd = new MySqlCommand(updateLastLoginTime, connection);
                    updateCmd.ExecuteNonQuery();
                    isUpdated = true;
                }
            }
            catch (Exception)
            {
                throw;
            }
            finally
            {
                if (connection.State == System.Data.ConnectionState.Open)
                {
                    connection.Close();
                }
            }
            return Json(isUpdated);
        }


        /*
        [CustomActionFilter]
        [OutputCache(NoStore = true, Duration = 0, VaryByParam = "*")]
        public ActionResult CreateAssetMapView()
        {
            ViewBag.Current = "Map View";
            return View();
        }
        */


        [CustomActionFilter]
        [OutputCache(NoStore = true, Duration = 0, VaryByParam = "*")]
        public ActionResult CreateAssetFormView()
        {
            ViewBag.Current = "Form View";
            return View();
        }
    }
}

