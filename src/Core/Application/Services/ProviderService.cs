using System.Threading.Tasks;
using Cyclo.Core.Application.DTOs;
using Cyclo.Core.Domain.Entities;
using Cyclo.Core.Domain.Repositories;

namespace Cyclo.Core.Application.Services
{
    public class ProviderService : IProviderService
    {
        private readonly IProviderRepository _providerRepository;
        private readonly ICreateCompanyService _createCompanyService;

        public ProviderService(
            IProviderRepository providerRepository,
            ICreateCompanyService createCompanyService)
        {
            _providerRepository = providerRepository;
            _createCompanyService = createCompanyService;
        }

        public async Task<Provider> CreateProviderAsync(CreateProviderDto dto)
        {
            Company company = null;
            
            if (dto.CreateCompany)
            {
                var companyDto = new CreateCompanyDto
                {
                    Name = dto.CompanyName,
                    Document = dto.CompanyDocument
                };
                
                company = await _createCompanyService.ExecuteAsync(companyDto);
            }

            var provider = new Provider(dto.Name, company?.Id);
            return await _providerRepository.CreateAsync(provider);
        }
    }

    public interface IProviderService
    {
        Task<Provider> CreateProviderAsync(CreateProviderDto dto);
    }
} 