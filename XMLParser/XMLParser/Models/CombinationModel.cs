using System;
using System.Collections.Generic;
using System.Linq;
using System.Web;
using System.Xml.Serialization;
using Newtonsoft.Json;

namespace XMLParser.Models
{
    public class CombinationModel
    {
        //[XmlAttribute("id")]
        public int id { get; set; }
        public List<SelectionQuestionModel> Selected { get; set; }
        public string NextQuestion { get; set; }
        public string Credits { get; set; }
    }
}