using System;
using System.Collections.Generic;
using System.Linq;
using System.Web;

namespace XMLParser.Models
{
    public class AnswersData
    {
        //[XmlAttribute("id")]
        public int id { get; set; }
        public string Answer { get; set; }
        public string File_name { get; set; } // question type
        public string Latitude { get; set; }  // questions latitude
        public string Longitude { get; set; }  // questions longitude
        public string Question { get; set; }
       
    }
}