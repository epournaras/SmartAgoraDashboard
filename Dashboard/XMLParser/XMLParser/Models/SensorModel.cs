using System;
using System.Collections.Generic;
using System.Linq;
using System.Web;
using System.Xml.Serialization;
using Newtonsoft.Json;

namespace XMLParser.Models
{
    public class SensorModel
    {
        //[XmlAttribute("id")]
        public int id { get; set; }
        public string Name { get; set; }
    }
}