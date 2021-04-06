using System;
using System.Collections.Generic;
using System.Linq;
using System.Web;
using System.Xml.Serialization;

namespace XMLParser.Models
{
    public class SelectionQuestionModel
    {
        //[XmlAttribute("id")]
        public int id { get; set; }
        public string Order { get; set; }
    }
}