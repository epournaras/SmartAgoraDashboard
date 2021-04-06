using System;
using System.Collections.Generic;
using System.Net;
using System.Net.Http;
using System.Net.Http.Headers;
using System.Threading.Tasks;
using System.Web.Script.Serialization;
using System.Text;
using Newtonsoft.Json;
using System.Diagnostics;
using System.Configuration;
using System.Reflection;
using System.IO;
using Newtonsoft.Json.Linq;
using HiveServer;
using System.Globalization;
using Nest;
using System.Web.Services;

namespace HiveServer
{
    public class Client
    {

        HttpClient client;
        private string host;
        private string port;
        private string url;

        public Client()
        {
            host = ConfigurationManager.AppSettings["host"];
            port = ConfigurationManager.AppSettings["port"];
            //DeleteProjectData("25042019180311007-dias");
            client = new HttpClient();
            //client.BaseAddress = new Uri("http://localhost:8080/");
            url = "http://" + host + ":" + port + "/";
            client.BaseAddress = new Uri("http://" + host + ":" + port + "/");
            client.DefaultRequestHeaders.Accept.Clear();
            client.DefaultRequestHeaders.Accept.Add(new MediaTypeWithQualityHeaderValue("application/json"));

        }

        #region Final Code
        public void CreateAsset(string projectId, Models.Asset[] Asset, string userName)
        {

            //try
            //{
            //    HttpResponseMessage response = await client.PostAsJsonAsync("/admin/projects/" + projectId + "/assets", Asset);
            //    response.EnsureSuccessStatusCode();
            //    Models.Asset asset = null;
            //    if (response.IsSuccessStatusCode)
            //    {
            //        var str = await response.Content.ReadAsStringAsync();
            //        JToken token = JObject.Parse(str);
            //        asset = ((Models.Asset)JsonConvert.DeserializeObject<Models.Asset>(token.First.First.ToString()));

            //    }
            //    return asset;
            //}
            //catch (Exception ex)
            //{
            //    throw ex;
            //}

            //----------------------------------------
            //--------For Windows, set false  --------
            //--------For Linux, set true     --------
            //----------------------------------------
            bool isLinuxDeployment = true;

            UriBuilder uri = new UriBuilder(Assembly.GetExecutingAssembly().CodeBase);
            string path = Path.GetDirectoryName(Uri.UnescapeDataString(uri.Path));

            //For Linux only
            if (isLinuxDeployment)
                path = path.Replace('\\', '/');
            //string fileName = DateTime.Now.ToString("ddMMyyyy_HHmmss", CultureInfo.InvariantCulture);
            string fileName = DateTime.Now.ToString("dd-MM-yyyy_HH-mm-ss-fff", CultureInfo.InvariantCulture);
            var json = JsonConvert.SerializeObject(Asset);
            try
            {
                // To create a string that specifies the path to a subfolder under your 
                // top-level folder, add a name for the subfolder to folderName.

                string pathString = System.IO.Path.Combine(path, userName);

                System.IO.Directory.CreateDirectory(pathString);

                if (isLinuxDeployment)
                    System.IO.File.WriteAllText(pathString + "/" + fileName + ".json", json);
                //System.IO.File.WriteAllText(path + "/asset.json", json);

                else
                    System.IO.File.WriteAllText(pathString + "\\" + fileName + ".json", json);
                //System.IO.File.WriteAllText(path + "\\asset.json", json);

                // Use ProcessStartInfo class
                ProcessStartInfo startInfo = new ProcessStartInfo();
                startInfo.CreateNoWindow = false;
                startInfo.UseShellExecute = false;
                startInfo.WindowStyle = ProcessWindowStyle.Normal;

                if (isLinuxDeployment)
                {
                    //startInfo.FileName = "/usr/bin/sudo";
                    //For Server Linux Deployment
                    //startInfo.Arguments = "/home/atif/Go_Workspace/bin/hive-command -host " + host + " -port " + port + " admin-create-assets " + projectId + " \"/" + path + "\"/asset.json";
                    //startInfo.Arguments = "/home/atif/Go_Workspace/bin/hive-command -host " + host + " -port " + port + " admin-create-assets " + projectId + " \"/" + pathString + "\"/" + fileName + ".json";

                    //For Ubuntu
                    startInfo.FileName = "/bin/bash";
                    startInfo.Arguments = "-c \"/home/atif/Go_Workspace/bin/hive-command -host " + host + " -port " + port + " admin-create-assets " + projectId + " \"/" + pathString + "\"/" + fileName + ".json\"";
                    //startInfo.Arguments = "-c \"/home/zaheer/Desktop/panecia/Go_Workspace/bin/hive-command -host " + host + " -port " + port + " admin-create-assets " + projectId + " \"/" + pathString + "\"/" + fileName + ".json\"";
                    //forhttp://mobileexperiment.inn.ac
                    //startInfo.FileName = "/bin/bash";
                    //startInfo.Arguments = "-c \"/home/mobexp/Dashboard/XMLParser/XMLParser/bin/hive-command -host " + host + " -port " + port + " admin-create-assets " + projectId + " \"/" + pathString + "\"/" + fileName + ".json\"";


                    //For Local linux deployment
                    //startInfo.Arguments = path + "/hive-command -host " + host + " -port " + port + " admin-create-assets basic \"/" + path + "\"/asset.json";
                }
                else
                {
                    startInfo.FileName = path + "\\hive-command.exe";
                    //startInfo.Arguments = "-host " + host + " -port " + port + " admin-create-assets " + projectId + " \"" + path + "\\asset.json\"";
                    startInfo.Arguments = "-host " + host + " -port " + port + " admin-create-assets " + projectId + " \"" + pathString + "\\" + fileName + ".json\"";

                }

                // Start the process with the info we specified.
                // Call WaitForExit and then the using statement will close.
                using (Process exeProcess = Process.Start(startInfo))
                {
                    exeProcess.WaitForExit();
                }

            }
            catch (Exception ex)
            {
                throw ex;
            }
        }

        public async Task<List<Models.Project>> GetAllProjects()
        {

            //HttpResponseMessage response = await client.GetAsync("http://localhost:8080/admin/projects");
            //string str = "";
            //if (response.IsSuccessStatusCode)
            //{
            //    str = await response.Content.ReadAsStringAsync();
            //}
            //return str;


            List<Models.Project> projects = null;
            //HttpResponseMessage response = await client.GetAsync("http://localhost:8080/admin/projects");
            //HttpResponseMessage response = await client.GetAsync(url + "admin/projects");
            HttpResponseMessage response = await client.GetAsync(url + "admin/projects?from=0&size=5000");

            if (response.IsSuccessStatusCode)
            {
                string str = await response.Content.ReadAsStringAsync();
                JToken token = JObject.Parse(str);

                projects = ((List<Models.Project>)JsonConvert.DeserializeObject<List<Models.Project>>(token.First.First.ToString()));
            }
            return projects;
        }

        public async Task<List<Models.Task>> GetTasks(string projectId)
        {
            List<Models.Task> task = null;
            HttpResponseMessage response = await client.GetAsync("/admin/projects/" + projectId + "/tasks?from=0&size=5000");

            if (response.IsSuccessStatusCode)
            {
                string str = await response.Content.ReadAsStringAsync();
                JToken token = JObject.Parse(str);
                task = ((List<Models.Task>)JsonConvert.DeserializeObject<List<Models.Task>>(token.First.First.ToString()));
            }
            return task;
        }

        public async Task<string> GetTasksData(string projectId, string taskId)
        {
            string str = "";
            HttpResponseMessage response = await client.GetAsync("/admin/projects/" + projectId + "/tasks/" + taskId);

            if (response.IsSuccessStatusCode)
            {
                str = await response.Content.ReadAsStringAsync();
            }
            return str;
        }
        public async Task<List<Models.Asset>> GetAssetAsync(string projectId)
        {
            List<Models.Asset> assets = null;
            HttpResponseMessage response = await client.GetAsync("/admin/projects/" + projectId + "/assets?from=0&size=5000");

            if (response.IsSuccessStatusCode)
            {
                string str = await response.Content.ReadAsStringAsync();
                JToken token = JObject.Parse(str);
                assets = ((List<Models.Asset>)JsonConvert.DeserializeObject<List<Models.Asset>>(token.First.First.ToString()));
            }
            return assets;
        }

        public async Task<Models.Asset> GetAssetData(string projectId, string assetId)
        {
            Models.Asset asset = null;
            HttpResponseMessage response = await client.GetAsync("/projects/" + projectId + "/assets/" + assetId);
            if (response.IsSuccessStatusCode)
            {
                string str = await response.Content.ReadAsStringAsync();
                JToken token = JObject.Parse(str);
                asset = ((Models.Asset)JsonConvert.DeserializeObject<Models.Asset>(token.First.First.ToString()));
            }
            return asset;
        }

        public static void DeleteProjectData(string projectId)
        {

            // ElasticClient elasticClient = null;

            //var uri = new Uri("http://localhost:9200");
            //var settings = new ConnectionSettings(uri);
            // settings.DefaultIndex("/hive/projects");
            // elasticClient = new ElasticClient(settings);

            //var response =  elasticClient.DeleteByQuery<string>(q => q
            // .Query(rq => rq
            // .Term(f => "_id", projectId)
            // )
            //  );

            //var uri = new Uri("http://localhost:9200/hive/projects/"+projectId+"/");
            //Console.WriteLine(uri.ToString());
            //var settings = new ConnectionSettings(uri);
            //settings.DefaultIndex("/hive/projects");
            //var elastic_client = new ElasticClient(settings);
            //var response = client.ClusterHealth();

          //  var request = new DeleteByQueryRequest<string>
           // {
             //   Query = new QueryContainer(
               // new TermQuery
            //{
              //  Field = "Id",
               // Value = projectId
            //}),
              //  Routing = "nest"
            //};
            // ExecuteAsync(projectId);
            // elastic_client.DeleteByQuery(request);
            //var results = elastic_client.Search<string>();
            //Console.WriteLine(results);
        }

        public void DeleteProject(string projectId)
        {
            var uri = new Uri("http://localhost:9200/");
            Console.WriteLine(uri.ToString());
            var settings = new ConnectionSettings(uri);
            //settings.DefaultIndex("/hive/projects");
            var elastic_client = new ElasticClient(settings);

            //var result = elastic_client.DeleteAsync(
            //DocumentPath<string>
            //        .Id(projectId),
            //            u => u
            //            .Index("hive")
            //           .Type("projects")
            //         );
            var result = elastic_client.Delete(new DeleteRequest("hive","projects",projectId));
            Console.WriteLine("Results: " + result);
        }

        public async Task<List<Models.User>> GetUserAsync(string projectId)
        {
            List<Models.User> users = null;
            HttpResponseMessage response = await client.GetAsync("/admin/projects/" + projectId + "/users?from=0&size=5000");

            if (response.IsSuccessStatusCode)
            {
                string str = await response.Content.ReadAsStringAsync();
                JToken token = JObject.Parse(str);
                users = ((List<Models.User>)JsonConvert.DeserializeObject<List<Models.User>>(token.First.First.ToString()));
            }
            return users;
        }

        public async Task<Models.User> GetUserData(string projectId, string userId)
        {
            Models.User user = null;
            HttpResponseMessage response = await client.GetAsync("/admin/projects/" + projectId + "/users/" + userId);

            if (response.IsSuccessStatusCode)
            {
                string str = await response.Content.ReadAsStringAsync();
                JToken token = JObject.Parse(str);
                user = ((Models.User)JsonConvert.DeserializeObject<Models.User>(token.ToString()));
            }
            return user;
        }

        public async Task<Models.Assignment> GetAssignmentData(string projectId, string assignmentId)
        {
            string state = "";
            Models.Assignment assignment = null;
            HttpResponseMessage response = await client.GetAsync("/projects/" + projectId + "/assignments/" + assignmentId);

            if (response.IsSuccessStatusCode)
            {
               var str = await response.Content.ReadAsStringAsync();

                JToken token = JObject.Parse(str);
                //state = token.Value<string>("name");
                assignment = ((Models.Assignment)JsonConvert.DeserializeObject<Models.Assignment>(token.First.First.ToString()));
            }
            return assignment;
        }


        public async Task<Models.Assignment> CreateAssignment(string projectId, string taskId, string assetId, string userId)
        {
            Models.Assignment assignment = null;
            //var baseAddress = new Uri("http://localhost:8080");
            var baseAddress = new Uri("http://" + host + ":" + port);
            var cookieContainer = new CookieContainer();
            using (var handler = new HttpClientHandler() { CookieContainer = cookieContainer })
            using (var client = new HttpClient(handler) { BaseAddress = baseAddress })
            {
                cookieContainer.Add(baseAddress, new Cookie(projectId + "_user_id", userId));

                HttpResponseMessage response = await client.GetAsync("/projects/" + projectId + "/tasks/" + taskId + "/assets/" + assetId + "/assignments");
                if (response.IsSuccessStatusCode)
                {
                    var str = await response.Content.ReadAsStringAsync();
                    JToken token = JObject.Parse(str);
                    assignment = ((Models.Assignment)JsonConvert.DeserializeObject<Models.Assignment>(token.ToString()));
                }
            }
            return assignment;
        }

        public async Task<Models.Project> CreateProject(Models.Project project)
        {
            try
            {
                HttpResponseMessage response = await client.PostAsJsonAsync("/admin/projects/" + project.Id, project);
                response.EnsureSuccessStatusCode();
                if (response.IsSuccessStatusCode)
                {
                    var str = await response.Content.ReadAsStringAsync();
                    JToken token = JObject.Parse(str);
                    project = ((Models.Project)JsonConvert.DeserializeObject<Models.Project>(token.First.First.ToString()));

                }
                return project;
            }
            catch (Exception ex)
            {
                throw ex;
            }
        }

        public async Task<Models.Task> CreateTask(string projectId, Models.Task task)
        {
            try
            {
                HttpResponseMessage response = await client.PostAsJsonAsync("/admin/projects/" + projectId + "/tasks/" + projectId + "-" + task.Name, task);
                response.EnsureSuccessStatusCode();
                if (response.IsSuccessStatusCode)
                {
                    var str = await response.Content.ReadAsStringAsync();
                    JToken token = JObject.Parse(str);
                    task = ((Models.Task)JsonConvert.DeserializeObject<Models.Task>(token.First.First.ToString()));
                }
                return task;
            }
            catch (Exception ex)
            {
                throw ex;
            }
        }

        #endregion

    }

}
