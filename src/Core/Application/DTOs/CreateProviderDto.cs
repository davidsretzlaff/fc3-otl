namespace Cyclo.Core.Application.DTOs
{
    public class CreateProviderDto
    {
        public string Name { get; set; }
        public bool CreateCompany { get; set; }
        public string CompanyName { get; set; }
        public string CompanyDocument { get; set; }
    }
} 