using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.IO;
using System.Linq;
using System.Net;
using System.Net.Http;
using System.Web.Http;
using System.Web;

namespace XMLParser.Controllers
{
    [RoutePrefix("api/xml")]
    public class XMLController : ApiController
    {
        [HttpPost]
        public IHttpActionResult HiveCall(string str)
        {
            HttpResponseMessage result = null;
            try
            {
                TextWriter writer = null;

                //var currentDate = DateTime.Now;
                ////var path = Directory.GetCurrentDirectory();
                ////var fullpath = path.Replace("\\", "/");

                //XMLFilePath = System.Web.Hosting.HostingEnvironment.MapPath("~/XMLFile" + currentDate.ToFileTimeUtc() + ".xml");
                //var serializer = new XmlSerializer(typeof(Questions));
                //XmlSerializerNamespaces ns = new XmlSerializerNamespaces();
                //ns.Add("", "");
                //writer = new StreamWriter(XMLFilePath);
                //serializer.Serialize(writer, mainModel, ns);
                if (writer != null)
                    writer.Close();
                result = Request.CreateResponse(HttpStatusCode.OK);
                return Ok();
            }
            catch (Exception ex)
            {
                Trace.TraceError(ex.InnerException == null ? ex.Message + " - " + ex.StackTrace : ex.Message + " - " + ex.InnerException.Message);
                throw ex;
            }
        }

    }
}
