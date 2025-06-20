﻿using Customer.Application.UseCases.User.Common;
using Customer.Application.UseCases.User.CreateUser;
using MediatR;

namespace Customer.Application.UseCases.Customer.CreateUser
{
    public interface ICreateCustomer : IRequestHandler<CreateCustomerInput, CustomerOutput>
    {
    }
}
