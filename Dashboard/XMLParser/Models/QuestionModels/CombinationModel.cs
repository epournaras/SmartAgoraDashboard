using System.Collections.Generic;

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