using System.Threading.Tasks;
using Cyclo.Core.Application.DTOs;
using Cyclo.Core.Domain.Entities;
using Cyclo.Core.Domain.Repositories;

namespace Cyclo.Core.Application.Services
{
    public class CreateCompanyService : ICreateCompanyService
    {
        private readonly ICompanyRepository _companyRepository;

        public CreateCompanyService(ICompanyRepository companyRepository)
        {
            _companyRepository = companyRepository;
        }

        public async Task<Company> ExecuteAsync(CreateCompanyDto dto)
        {
            var company = new Company(dto.Name, dto.Document);
            return await _companyRepository.CreateAsync(company);
        }
    }

    public interface ICreateCompanyService
    {
        Task<Company> ExecuteAsync(CreateCompanyDto dto);
    }
} 