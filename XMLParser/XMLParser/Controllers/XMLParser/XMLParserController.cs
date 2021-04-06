using System;
using System.Diagnostics;
using System.Collections.Generic;
using System.Linq;
using System.Net;
using System.Net.Http;
using System.Web.Http;
using XMLParser.XMLParserInterfaces;
using XMLParser.Models;
using System.Xml.Serialization;
using System.Text;
using System.IO;
using System.Web;
using HiveServer;
using Newtonsoft.Json;
using System.Globalization;
using Models.HiveModels;
using MySql.Data.MySqlClient;
using HiveServer.Models;
using Models.QuestionModels;
using System.Configuration;

namespace XMLParser.Controllers.XMLParser
{
    [RoutePrefix("api/xmlparser")]
    public class XMLParserController : ApiController
    {
        Client client = new Client();
        static string ConnectionString = ConfigurationManager.AppSettings["connectionString"];
        #region private properties
        //private readonly IXMLParser xmlparserService;

        public static string JSONFilePath;
        #endregion

        [HttpPost]
        [Route("GenerateHiveCall")]
        public void GenerateHiveCall(ProjectQuestionModel projectQuestionModel)
        {
            try
            {

                string fileName = projectQuestionModel.QuestionsModel.StartAndDestinationModel.Mode + "_" + DateTime.Now.ToString("ddMMyyyy_HHmmss", CultureInfo.InvariantCulture);
                HiveServer.Models.Asset[] asset = new HiveServer.Models.Asset[] {
                    new HiveServer.Models.Asset
                    {
                        Name = fileName,
                        Url = "smart-agora.org",
                        Metadata = new HiveServer.Models.Metadata()
                        {
                            record=projectQuestionModel.QuestionsModel
                        }
                    }
                };

                //string email = System.Web.HttpContext.Current.Session["Email"].ToString();
                //string userName = email.Substring(0, email.IndexOf("@"));//"userName"; // System.Web.HttpContext.Current.Session["Email"].ToString();
                string userName = "userName";
                client.CreateAsset(projectQuestionModel.ProjectId, asset, userName);
            }
            catch (Exception ex)
            {
                Trace.TraceError(ex.InnerException == null ? ex.Message + " - " + ex.StackTrace : ex.Message + " - " + ex.InnerException.Message);
                throw ex;
            }
        }

        [HttpPost]
        [Route("GenerateXMLFile")]
        public IHttpActionResult GenerateXMLFile(ProjectQuestionModel mainModel)
        {
            try
            {
                var json = JsonConvert.SerializeObject(mainModel, Formatting.Indented, new JsonSerializerSettings { DefaultValueHandling = DefaultValueHandling.Ignore });

                var currentDate = DateTime.Now;
                //var path = Directory.GetCurrentDirectory();
                //var fullpath = path.Replace("\\", "/");

                JSONFilePath = System.Web.Hosting.HostingEnvironment.MapPath("~/JSONFile" + currentDate.ToFileTimeUtc() + ".json");
                System.IO.File.WriteAllText(JSONFilePath, json);
                //TextWriter writer = null;

                //var currentDate = DateTime.Now;
                ////var path = Directory.GetCurrentDirectory();
                ////var fullpath = path.Replace("\\", "/");

                //XMLFilePath = System.Web.Hosting.HostingEnvironment.MapPath("~/XMLFile" + currentDate.ToFileTimeUtc() + ".xml");
                //var serializer = new XmlSerializer(typeof(Questions));
                //XmlSerializerNamespaces ns = new XmlSerializerNamespaces();
                //ns.Add("", "");
                //writer = new StreamWriter(XMLFilePath);
                //serializer.Serialize(writer, mainModel, ns);
                //if (writer != null)
                //    writer.Close();
                return Ok();
            }
            catch (Exception ex)
            {
                Trace.TraceError(ex.InnerException == null ? ex.Message + " - " + ex.StackTrace : ex.Message + " - " + ex.InnerException.Message);
                throw ex;
            }
        }

        [HttpGet]
        public HttpResponseMessage GetDownloadFile()
        {
            HttpResponseMessage result = null;
            //var localFilePath = HttpContext.Current.Server.MapPath(@"E:/XMLFiles/XMLFile131416793102014093.xml");
            var localFilePath = JSONFilePath;
            if (!File.Exists(localFilePath))
            {
                result = Request.CreateResponse(HttpStatusCode.Gone);
            }
            else
            {
                // Serve the file to the client
                result = Request.CreateResponse(HttpStatusCode.OK);
                result.Content = new StreamContent(new FileStream(localFilePath, FileMode.Open, FileAccess.Read));
                result.Content.Headers.ContentDisposition = new System.Net.Http.Headers.ContentDispositionHeaderValue("attachment");
                result.Content.Headers.ContentDisposition.FileName = "JSONFile";
            }
            return result;
        }

        [HttpDelete]
        public IHttpActionResult DeleteDownloadedFile()
        {
            try
            {
                var localFilePath = JSONFilePath;
                if (File.Exists(localFilePath))
                {
                    File.Delete(localFilePath);
                }
                return Ok(true);
            }
            catch (Exception ex)
            {
                Trace.TraceError(ex.InnerException == null ? ex.Message + " - " + ex.StackTrace : ex.Message + " - " + ex.InnerException.Message);
                throw ex;
            }
        }

        [HttpGet]
        [Route("GetAllProjects")]
        public async System.Threading.Tasks.Task<List<HiveServer.Models.Project>> GetAllProjects(Questions mainModel)
        {
            return await client.GetAllProjects();
        }

        //[HttpPost]
        //[Route("CreateUserAnswer")]
        //public IHttpActionResult CreateUserAnswer(AnswersData mainModel)
        //{
        //    MySqlConnection connection = new MySqlConnection(ConnectionString);
        //    bool validModel = false;
        //    try
        //    {
        //        connection.Open();
        //        string query = "INSERT INTO answersdata(answer, files__name, latitude, longitude, question) VALUES ('" +
        //            mainModel.Answer + "','" +
        //            mainModel.File_name + "','" +
        //            mainModel.Latitude + "','" +
        //            mainModel.Longitude + "','" +
        //            mainModel.Question + "')";
        //        MySqlCommand cmd = new MySqlCommand(query, connection);
        //        int rows = cmd.ExecuteNonQuery();
        //        if (rows > 0)
        //        {
        //            validModel = true;
        //            return Ok();
        //        }
        //        return BadRequest();
        //    }
        //    catch (Exception ex)
        //    {
        //        Trace.TraceError(ex.InnerException == null ? ex.Message + " - " + ex.StackTrace : ex.Message + " - " + ex.InnerException.Message);
        //        throw ex;
        //    }
        //}

        [HttpPost]
        [Route("SaveAssignment")]
        public IHttpActionResult SaveAssignment(Assignment assignmentModel)
        {
            MySqlConnection connection = new MySqlConnection(ConnectionString);
            int SED_id = 0;
            int assignmentrows = 0;
            bool validModel = false;
            try
            {
                connection.Open();
                if (assignmentModel != null)
                {
                    string assignmentquery = "INSERT INTO assignment(Id, User, Project, Task, State) VALUES ('" +
                    assignmentModel.Id + "','" +
                    assignmentModel.User + "','" +
                    assignmentModel.Project + "','" +
                    assignmentModel.Task + "','" +
                    assignmentModel.State + "')";

                    MySqlCommand assignmentcmd = new MySqlCommand(assignmentquery, connection);
                    assignmentrows = assignmentcmd.ExecuteNonQuery();
                    //int id = Convert.ToInt32(cmd.LastInsertedId);
                    if (assignmentModel.Asset.Metadata.record.StartAndDestinationModel != null)
                    {
                        var startEndDestinationObj = assignmentModel.Asset.Metadata.record.StartAndDestinationModel;
                        string startEndDestinationquery = "INSERT INTO startanddestination(DefaultCredit, DestinationLatitude, DestinationLongitude, Mode, StartLatitude, StartLongitude) VALUES ('" +
                        startEndDestinationObj.DefaultCredit + "','" +
                        startEndDestinationObj.DestinationLatitude + "','" +
                        startEndDestinationObj.DestinationLongitude + "','" +
                        startEndDestinationObj.Mode + "','" +
                        startEndDestinationObj.StartLatitude + "','" +
                        startEndDestinationObj.StartLongitude + "')";

                        MySqlCommand SEDcmd = new MySqlCommand(startEndDestinationquery, connection);
                        SEDcmd.ExecuteNonQuery();
                        SED_id = Convert.ToInt32(SEDcmd.LastInsertedId);
                    }

                    if (SED_id > 0 && assignmentModel.Asset != null)
                    {
                        var assetObj = assignmentModel.Asset;
                        string assetquery = "INSERT INTO asset(Id, AssignmentId, Name, StartEndDestinationId, Url) VALUES ('" +
                        assetObj.Id + "','" +
                        assignmentModel.Id + "','" +
                        assetObj.Name + "','" +
                        SED_id + "','" +
                        assetObj.Url + "')";
                        MySqlCommand assetcmd = new MySqlCommand(assetquery, connection);
                        assetcmd.ExecuteNonQuery();
                    }

                    if (assignmentModel.Asset.Metadata.record.SampleDataModel != null && assignmentModel.Asset.Metadata.record.SampleDataModel.Any())
                    {
                        var lstSampleDateModel = assignmentModel.Asset.Metadata.record.SampleDataModel;
                        for (var i = 0; i < lstSampleDateModel.Count; i++)
                        {
                            var questionData_Id = 0;
                            var questionDataObj = lstSampleDateModel[i];
                            string questionDataquery = "INSERT INTO questiondata(QuestionId, AssetId, Frequency, Latitude, Longitude, Mandatory, Question, Sequence, Time, Type, Visibility,Vicinity) VALUES ('" +
                            questionDataObj.id + "','" +
                            assignmentModel.Asset.Id + "','" +
                            questionDataObj.Frequency + "','" +
                            questionDataObj.Latitude + "','" +
                            questionDataObj.Longitude + "','" +
                            questionDataObj.Mandatory + "','" +
                            questionDataObj.Question + "','" +
                            questionDataObj.Sequence + "','" +
                            questionDataObj.Time + "','" +
                            questionDataObj.Type + "','" +
                            questionDataObj.Visibility + "','" +
                            questionDataObj.Vicinity + "')";
                            MySqlCommand questionDatacmd = new MySqlCommand(questionDataquery, connection);
                            questionDatacmd.ExecuteNonQuery();
                            //questionData_Id = Convert.ToInt32(questionDatacmd.LastInsertedId);
                            if (lstSampleDateModel[i].Option != null && lstSampleDateModel[i].Option.Any())
                            {
                                for (var optionI = 0; optionI < lstSampleDateModel[i].Option.Count; optionI++)
                                {
                                    var questionOptionObj = lstSampleDateModel[i].Option[optionI];
                                    string questionOptionquery = "INSERT INTO questionoption(QuestionOptionId, AssignmentId, AssetId, Name, Credits, NextQuestion, QuestionDataId) VALUES ('" +
                                    questionOptionObj.id + "','" +
                                    assignmentModel.Id + "','" +
                                    assignmentModel.Asset.Id + "','" +
                                    questionOptionObj.Name + "','" +
                                    questionOptionObj.Credits + "','" +
                                    questionOptionObj.NextQuestion + "','" +
                                    questionDataObj.id + "')";
                                    MySqlCommand questionOptioncmd = new MySqlCommand(questionOptionquery, connection);
                                    questionOptioncmd.ExecuteNonQuery();
                                }
                            }
                            if (lstSampleDateModel[i].Sensor != null && lstSampleDateModel[i].Sensor.Any())
                            {
                                for (var sensorI = 0; sensorI < lstSampleDateModel[i].Sensor.Count; sensorI++)
                                {
                                    var sensorObj = lstSampleDateModel[i].Sensor[sensorI];
                                    string questionsensorquery = "INSERT INTO questionsensor(SensorId, QuestionDataId) VALUES ('" +
                                    sensorObj.id + "','" +
                                    questionDataObj.id + "')";
                                    MySqlCommand sensorcmd = new MySqlCommand(questionsensorquery, connection);
                                    sensorcmd.ExecuteNonQuery();
                                }
                            }
                            if (lstSampleDateModel[i].Combination != null && lstSampleDateModel[i].Combination.Any())
                            {
                                for (var combinationI = 0; combinationI < lstSampleDateModel[i].Combination.Count; combinationI++)
                                {
                                    /* perform action */
                                }
                            }
                        }
                    }

                    if (assignmentModel.SubmittedAnswerData != null && assignmentModel.SubmittedAnswerData.SubmittedData.Any())
                    {
                        for (int i = 0; i < assignmentModel.SubmittedAnswerData.SubmittedData.Count; i++)
                        {
                            string queryanswerdata = "INSERT INTO answersdata(Answer, Files_Name, Latitude, Longitude, Question, Type, AssignmentId, AssetId,TimeAtAnswering) VALUES ('" +
                                                       assignmentModel.SubmittedAnswerData.SubmittedData[i].Answer + "','" +
                                                       assignmentModel.SubmittedAnswerData.SubmittedData[i].Files_Name + "','" +
                                                       assignmentModel.SubmittedAnswerData.SubmittedData[i].Latitude + "','" +
                                                       assignmentModel.SubmittedAnswerData.SubmittedData[i].Longitude + "','" +
                                                       assignmentModel.SubmittedAnswerData.SubmittedData[i].Question + "','" +
                                                       assignmentModel.SubmittedAnswerData.SubmittedData[i].Type + "','" +
                                                       assignmentModel.Id + "','" +
                                                       assignmentModel.Asset.Id + "','" +
                                                       assignmentModel.SubmittedAnswerData.SubmittedData[i].TimeAtAnswering + "')";
                            MySqlCommand answerdatacmd = new MySqlCommand(queryanswerdata, connection);
                            answerdatacmd.ExecuteNonQuery();
                        }

                    }


                    //MySqlCommand cmd = new MySqlCommand(query, connection);
                    //int rows = cmd.ExecuteNonQuery();
                    //int id = Convert.ToInt32(cmd.LastInsertedId);
                }
                if (assignmentrows > 0)
                {
                    validModel = true;
                    return Ok();
                }
                return BadRequest();
            }
            catch (Exception ex)
            {
                Trace.TraceError(ex.InnerException == null ? ex.Message + " - " + ex.StackTrace : ex.Message + " - " + ex.InnerException.Message);
                throw ex;
            }
        }

        [HttpPost]
        [Route("SaveSensorData")]
        public IHttpActionResult SaveSensorData(SensorDataModel sensorDataModel)
        {
            MySqlConnection connection = new MySqlConnection(ConnectionString);
            int sensorRows = 0;

            try
            {
                connection.Open();
                if (sensorDataModel != null)
                {

                    string query = "INSERT INTO sensordata(Acceleration, Frequency, Gyroscope, Light, Location, Noise, Proximity, Question, QuestionId, Id, AssignmentId,TimeAtSensoring) VALUES ('" +
                        sensorDataModel.Acceleration + "','" +
                        sensorDataModel.Frequency + "','" +
                        sensorDataModel.Gyroscope + "','" +
                        sensorDataModel.Light + "','" +
                        sensorDataModel.Location + "','" +
                        sensorDataModel.Noise + "','" +
                        sensorDataModel.Proximity + "','" +
                        sensorDataModel.Question + "','" +
                        sensorDataModel.QuestionId + "','" +
                        sensorDataModel.Id + "','" +
                        sensorDataModel.AssignmentId + "','" +
                        sensorDataModel.TimeAtSensoring + "')";

                    MySqlCommand sqlCommand = new MySqlCommand(query, connection);
                    sensorRows = sqlCommand.ExecuteNonQuery();
                }

                if (sensorRows > 0)
                {
                    return Ok();
                }
                return BadRequest();

            }
            catch (Exception ex)
            {
                Trace.TraceError(ex.InnerException == null ? ex.Message + " - " + ex.StackTrace : ex.Message + " - " + ex.InnerException.Message);
                throw ex;
            }
        }

    }
}
