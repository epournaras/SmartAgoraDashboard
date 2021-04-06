using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Linq;
using System.Web;
using System.Xml.Serialization;
using Newtonsoft.Json;

namespace XMLParser.Models
{
    public class QuestionDataModel
    {
        //[XmlAttribute("id")]
        public int id { get; set; }
        public string Question { get; set; }
        public string Type { get; set; } // question type
        public string Latitude { get; set; }  // questions latitude
        public string Longitude { get; set; }  // questions longitude
        public List<SensorModel> Sensor { get; set; }
        public string Time { get; set; }
        public string Frequency { get; set; }
        [DefaultValue("Disable")]
        public string Sequence { get; set; }
        public string Visibility { get; set; }
        public string Mandatory { get; set; }
        public List<OptionModel> Option { get; set; }
        [DefaultValue(null)]
        public List<CombinationModel> Combination { get; set; }
    }
}