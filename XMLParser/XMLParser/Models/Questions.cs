using System;
using System.Collections.Generic;
using System.Linq;
using System.Web;
using System.Xml.Serialization;
using Newtonsoft.Json;

namespace XMLParser.Models
{
    public class Questions
    {
        public StartAndDestinationModel StartAndDestinationModel { get; set; }
        public List<QuestionDataModel> SampleDataModel { get; set; }
    }
}