using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Models.QuestionModels
{
    public class SensorDataModel
    {
        public int Id { get; set; }
        public string Noise{ get; set; }
        public string Acceleration{ get; set; }
        public string Light{ get; set; }
        public string Gyroscope{ get; set; }
        public string Proximity{ get; set; }
        public string Location { get; set; }
        public string Frequency { get; set; }
        public string Question { get; set; }
        public string QuestionId { get; set; }
        public string AssignmentId { get; set; }
        public string TimeAtSensoring { get; set; }
    }
}
