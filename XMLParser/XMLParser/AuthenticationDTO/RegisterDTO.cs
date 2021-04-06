using System;
using System.Collections.Generic;
using System.Linq;
using System.Web;

namespace XMLParser.AuthenticationDTO
{
    public class RegisterDTO
    {
        public string FirstName { get; set; }
        public string LastName { get; set; }
        //public DateTime DOB { get; set; }
        //public string City { get; set; }
        //public string Status { get; set; }
        public string StudentId { get; set; }
        public string EmailId { get; set; }
        public string Password { get; set; }
        public DateTime registrationTime { get; set; }
        public DateTime LastLoginTime { get; set; }
        public string username { get; set; }

    }
}