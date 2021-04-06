using System.ComponentModel;
namespace XMLParser.Models
{
    public class TextBoxModel
    {
        public string NextQuestion { get; set; }
        [DefaultValue("Disable")]
        public string Credits { get; set; }
    }
}
