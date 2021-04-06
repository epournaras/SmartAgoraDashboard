using System.ComponentModel;

namespace XMLParser.Models
{
    public class OptionModel
    {
        //[XmlAttribute("id")]
        public int id { get; set; }
        public string Name { get; set; }
        [DefaultValue("Disable")]
        public string NextQuestion { get; set; }
        [DefaultValue("Disable")]
        public string Credits { get; set; }
    }
}