using System;
using System.Collections.Generic;
using System.Linq;
using System.Web;
using XMLParser.XMLParserInterfaces;
using XMLParser.Models;
using System.Xml.Serialization;
using System.Text;
using System.IO;

namespace XMLParser.XMLParserBLL
{
    public class XMLBLL : IXMLParser
    {
        public void GenerateXMLFile(Questions mainModel)
        {
            TextWriter writer = null;
            try
            {
                var currentDate = DateTime.Now;
                var XMLFilePath = "../XMLFiles" + currentDate.ToFileTimeUtc();
                var serializer = new XmlSerializer(typeof(Questions));
                writer = new StreamWriter(XMLFilePath);
                serializer.Serialize(writer, mainModel);
            }
            finally
            {
                if (writer != null)
                    writer.Close();
            }
        }
    }
}