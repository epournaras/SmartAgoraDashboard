using HiveServer;
using HiveServer.Models;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using System.Web.Mvc;
using XMLParser.Models;
using XMLParser.Filter;
using System.Web.Routing;
using System.Globalization;
using System.Web.UI;

namespace XMLParser.Controllers
{
    [CustomActionFilter]
    public class HiveController : Controller
    {
        // GET: Hive
        public ActionResult Index()
        {
            return View();
        }

        public async Task<ActionResult> GetAllProjects(jQueryDataTableParamModel param)
        {
            Client client = new Client();
            try
            {
                List<Project> projects = await client.GetAllProjects();
                var result = from c in projects
                             where c.Id.Contains("-") && ((c.Id.Substring(c.Id.LastIndexOf("-")).Equals("-" + Session["Email"].ToString().Substring(0, Session["Email"].ToString().IndexOf("@")))) || c.Id.Substring(c.Id.LastIndexOf("-")).Equals("-" + Session["UserName"].ToString()))
                             select new Project
                             {
                                 Id = c.Id,
                                 Name = c.Name,
                                 Description = c.Description.Length > 0 && c.Description.Contains("#")? c.Description.Substring(0,c.Description.IndexOf('#')): c.Description,
                                 AssetCount = c.AssetCount,
                                 TaskCount = c.TaskCount,
                                 UserCount = c.UserCount,
                                 AssignmentCount = c.AssignmentCount,
                                 MetaProperties = c.MetaProperties
                             };
                return Json(new
                {
                    sEcho = param.sEcho,
                    iTotalRecords = result.Count(),
                    iTotalDisplayRecords = result.Count(),
                    aaData = result
                }, JsonRequestBehavior.AllowGet);
            }
            catch (Exception ex)
            {
                throw ex;
            }
        }

        public async Task<ActionResult> GetProjectsIds()
        {
            HiveServer.Client client = new HiveServer.Client();
            try
            {
                List<Project> projects = await client.GetAllProjects();
                var result = from c in projects
                             where c.Id.Contains("-") && ((c.Id.Substring(c.Id.LastIndexOf("-")).Equals("-" + Session["Email"].ToString().Substring(0, Session["Email"].ToString().IndexOf("@")))) || c.Id.Substring(c.Id.LastIndexOf("-")).Equals("-" + Session["UserName"].ToString()))
                             select new Project
                             {
                                 Id = c.Id,
                                 Name = c.Name
                             };
                return Json(new
                {
                    iTotalRecords = result.Count(),
                    iTotalDisplayRecords = result.Count(),
                    aaData = result
                }, JsonRequestBehavior.AllowGet);
            }
            catch (Exception ex)
            {
                throw ex;
            }
        }

        public async Task<ActionResult> GetAssetData(string projectId, string assetId)
        {
            try
            {
                Client client = new Client();
                Asset assets = await client.GetAssetData(projectId, assetId);
                return Json(new
                {
                    aaData = assets
                }, JsonRequestBehavior.AllowGet);
            }
            catch (Exception ex)
            {
                throw ex;
            }
        }

        public async Task<ActionResult> GetAssetsIds(string projectId)
        {
            try
            {
                Client client = new Client();
                List<Asset> assets = await client.GetAssetAsync(projectId);
                if (assets != null)
                {
                    var result = from c in assets
                                 select new Asset
                                 {
                                     Id = c.Id,
                                     Name = c.Name,
                                 };

                    return Json(new
                    {
                        iTotalRecords = result.Count(),
                        iTotalDisplayRecords = result.Count(),
                        aaData = result
                    }, JsonRequestBehavior.AllowGet);
                }
                else
                {
                    return Json(new
                    {
                        iTotalRecords = 0,
                        iTotalDisplayRecords = 0,
                        assetFoundMessage = "noAssetfound"
                    }, JsonRequestBehavior.AllowGet);

                }
            }
            catch (Exception ex)
            {
                throw ex;
            }
        }

        public async Task<ActionResult> GetUserData(string projectId, string userId)
        {
            try
            {
                Client client = new Client();
                User user = await client.GetUserData(projectId, userId);

                return Json(new
                {
                    aaData = user
                }, JsonRequestBehavior.AllowGet);
            }
            catch (Exception ex)
            {
                throw ex;
            }
        }

        public async Task<ActionResult> GetUsersIds(string projectId)
        {
            try
            {
                Client client = new Client();
                List<User> users = await client.GetUserAsync(projectId);
                if (users != null)
                {
                    var result = from c in users
                                 select new User
                                 {
                                     Id = c.Id,
                                     Name = c.Name
                                 };

                    return Json(new
                    {
                        iTotalRecords = result.Count(),
                        iTotalDisplayRecords = result.Count(),
                        aaData = result
                    }, JsonRequestBehavior.AllowGet);
                }
                else
                {
                    return Json(new
                    {
                        //sEcho = param.sEcho
                        iTotalRecords = 0,
                        iTotalDisplayRecords = 0,
                        userFoundMessage = "nouserfound"
                    }, JsonRequestBehavior.AllowGet);
                }
            }
            catch (Exception ex)
            {
                throw ex;
            }
        }

        public async Task<ActionResult> GetTaskData(string projectId, string taskId)
        {
            try
            {
                Client client = new Client();
                string result = await client.GetTasksData(projectId, taskId);

                return Json(new
                {
                    iTotalRecords = 1,
                    iTotalDisplayRecords = 1,
                    aaData = result
                }, JsonRequestBehavior.AllowGet);
            }
            catch (Exception ex)
            {
                throw ex;
            }
        }

        public async Task<ActionResult> GetTasksIds(string projectId)
        {
            try
            {
                Client client = new Client();
                List<HiveServer.Models.Task> tasks = await client.GetTasks(projectId);
                if (tasks != null)
                {
                    var result = from c in tasks
                                 select new HiveServer.Models.Task
                                 {
                                     Id = c.Id,
                                     Name = c.Name,
                                 };

                    return Json(new
                    {
                        iTotalRecords = result.Count(),
                        iTotalDisplayRecords = result.Count(),
                        aaData = result
                    }, JsonRequestBehavior.AllowGet);
                }
                else
                {
                    return Json(new
                    {
                        //sEcho = param.sEcho
                        iTotalRecords = 0,
                        iTotalDisplayRecords = 0,
                        taskFoundMessage = "notaskfound"
                    }, JsonRequestBehavior.AllowGet);
                }
            }
            catch (Exception ex)
            {
                throw ex;
            }
        }

        public async Task<ActionResult> CreateAssignment(string projectId, string taskId, string assetId, string userIdsList)
        {
            try
            {
                var userIds = userIdsList.Split(',');
                Client client = new Client();
                List<Assignment> assignmentResults = new List<Assignment>();
                Assignment assignments = new Assignment();
               
                foreach (var userId in userIds)
                {
                    string user = userId.Trim(new Char[] { '[', '\\', ']', '"' });
                    String assignmentId = projectId + "HIVE" + taskId + "HIVE" + assetId + "HIVE" + user;
                    assignments = await client.GetAssignmentData(projectId, assignmentId);
                    String state = assignments.State;

                    if (state.Equals("unfinished")) {

                    }
                    else if(state.Equals("finished"))
                    {
                        Assignment result = await client.CreateAssignment(projectId, taskId, assetId, user);
                        assignmentResults.Add(result);
                    }
                    Console.WriteLine("STATE:" + state);
                }

                return Json(new
                {
                    aaData = assignmentResults
                }, JsonRequestBehavior.AllowGet);
            }
            catch (Exception ex)
            {
                throw ex;
            }
        }

        public void DeleteProject(string projectId)
        {
            try
            {
                Client client = new Client();
                client.DeleteProject(projectId);
                }
            catch (Exception ex)
            {
                throw ex;
            }
        }

        public async Task<ActionResult> CreateProject(string projectId, string projectName, string projectDesc)
        {
            try
            {
                if(projectId.Equals(""))
                    projectId = DateTime.Now.ToString("ddMMyyyyHHmmssfff", CultureInfo.InvariantCulture);
                String dataScientistId = "";


                if (Session["UserName"].ToString().Equals("")) {
                    if (Session["Email"] != null)
                    {
                        string email = Session["Email"].ToString();
                        dataScientistId = email.Substring(0, email.IndexOf('@'));
                    }
                }
                else
                {
                    dataScientistId = Session["UserName"].ToString();
                }
                Client client = new Client();
                Project project = new Project
                {
                    Id = projectId + "-" + dataScientistId,
                    Name = projectName,
                    Description = projectDesc
                };
                var result = await client.CreateProject(project);

                return Json(new
                {
                    aaData = result
                }, JsonRequestBehavior.AllowGet);
            }
            catch (Exception ex)
            {
                throw ex;
            }
        }

        public async Task<ActionResult> CreateTask(string projectId, string name, string desc, string state, int total=500, int matching=500)
        {
            try
            {
                Client client = new Client();
                HiveServer.Models.Task task = new HiveServer.Models.Task
                {
                    Name = name,
                    Description = desc,
                    CurrentState = state,
                    CompletionCriteria = new CompletionCriteria
                    {
                        Matching = matching,
                        Total = total
                    }
                };
                var result = await client.CreateTask(projectId, task);

                return Json(new
                {
                    aaData = result
                }, JsonRequestBehavior.AllowGet);
            }
            catch (Exception ex)
            {
                throw ex;
            }
        }
    }
}
