using System;
using System.Collections.Generic;
using System.Linq;
using System.Web;
using System.Xml.Serialization;

namespace XMLParser.Models
{
    public class StartAndDestinationModel
    {
        // user start latitude
        public string StartLatitude { get; set; }
        // user start longitude
        public string StartLongitude { get; set; }
        // user destination latitude
        public string DestinationLatitude { get; set; }
        // user destination longitude
        public string DestinationLongitude { get; set; }
        public string Mode { get; set; }
        public string DefaultCredit { get; set; }
    }
}