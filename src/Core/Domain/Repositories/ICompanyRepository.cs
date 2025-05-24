using System;
using System.Threading.Tasks;
using Cyclo.Core.Domain.Entities;

namespace Cyclo.Core.Domain.Repositories
{
    public interface ICompanyRepository
    {
        Task<Company> CreateAsync(Company company);
    }
} 