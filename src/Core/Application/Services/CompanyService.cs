using System.Threading.Tasks;
using Cyclo.Core.Application.DTOs;
using Cyclo.Core.Domain.Entities;
using Cyclo.Core.Domain.Repositories;

namespace Cyclo.Core.Application.Services
{
    public class CompanyService : ICompanyService
    {
        private readonly ICompanyRepository _companyRepository;

        public CompanyService(ICompanyRepository companyRepository)
        {
            _companyRepository = companyRepository;
        }

        public async Task<Company> CreateCompanyAsync(CreateCompanyDto dto)
        {
            var company = new Company(dto.Name, dto.Document);
            return await _companyRepository.CreateAsync(company);
        }
    }

    public interface ICompanyService
    {
        Task<Company> CreateCompanyAsync(CreateCompanyDto dto);
    }
} 