using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using XMLParser.Models;

namespace XMLParser.XMLParserInterfaces
{
    public interface IXMLParser
    {
        void GenerateXMLFile(Questions mainModel);
    }
}
