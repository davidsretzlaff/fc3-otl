namespace Customer.Domain.Domain
{
    public class Customer
    {
        public Guid Id { get; private set; }
        public string Name { get; private set; }
        public string Email {get; private set;}
        
        private Customer() { }
        
        public static Customer Create(string name, string email)
        {
            var customer = new Customer
            {
                Id = Guid.NewGuid(),
                Name = name,
                Email = email
            };
            return customer;
        }
    }
}
