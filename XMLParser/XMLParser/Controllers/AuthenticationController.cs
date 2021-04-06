using System;
using System.Collections.Generic;
using System.Configuration;
using System.Linq;
using System.Net;
using System.Net.Http;
using System.Web;
using System.Web.Http;
using MySql.Data.MySqlClient;
using XMLParser.AuthenticationDTO;

namespace XMLParser.Controllers
{
    [RoutePrefix("api/authentication")]
    public class AuthenticationController : ApiController
    {
        static string ConnectionString = ConfigurationManager.AppSettings["connectionString"];
        static string Username = null;
        

        
    }
}
