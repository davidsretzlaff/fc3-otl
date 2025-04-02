using Microsoft.Data.Sqlite;
using OpenTelemetry.Resources;
using OpenTelemetry.Trace;

var builder = WebApplication.CreateBuilder(args);

// Add services to the container.
builder.Services.AddControllers();

// Configure OpenTelemetry
builder.Services.AddOpenTelemetry()
    .WithTracing(tracerProviderBuilder =>
    {
        tracerProviderBuilder
            .SetResourceBuilder(ResourceBuilder.CreateDefault().AddService("MyService"))
            .AddAspNetCoreInstrumentation()
            .AddHttpClientInstrumentation()
            .AddSqlClientInstrumentation()
            .AddZipkinExporter(options =>
            {
                options.Endpoint = new Uri("http://localhost:9411/api/v2/spans");
            });
    });

var app = builder.Build();

// Initialize SQLite database
InitializeDatabase();

// Configure the HTTP request pipeline.
if (app.Environment.IsDevelopment())
{
    app.UseDeveloperExceptionPage();
}

app.UseRouting();

app.UseAuthorization();

app.MapControllers();

// Add a simple root endpoint
app.MapGet("/", () => "Hello World!");

app.Run("http://0.0.0.0:80");

void InitializeDatabase()
{
    using var connection = new SqliteConnection("Data Source=mydatabase.db");
    connection.Open();

    var command = connection.CreateCommand();
    command.CommandText = @"
        CREATE TABLE IF NOT EXISTS MyTable (
            Id INTEGER PRIMARY KEY AUTOINCREMENT,
            Name TEXT NOT NULL
        );
    ";
    command.ExecuteNonQuery();

    // Check if the table already contains data
    command.CommandText = "SELECT COUNT(*) FROM MyTable";
    var count = (long)command.ExecuteScalar();

    if (count == 0)
    {
        // Insert 10 initial records
        for (int i = 1; i <= 10; i++)
        {
            command.CommandText = $@"
                INSERT INTO MyTable (Name) VALUES ('Initial Data {i}');
            ";
            command.ExecuteNonQuery();
        }
    }
}

