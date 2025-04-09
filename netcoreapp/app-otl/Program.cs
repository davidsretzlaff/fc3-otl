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
            .SetResourceBuilder(ResourceBuilder.CreateDefault()
                .AddService("DotNetService"))
            .AddAspNetCoreInstrumentation()
            .AddHttpClientInstrumentation()
            .AddSqlClientInstrumentation()
            .AddSource("MyController")
            .AddOtlpExporter(options =>
            {
                options.Endpoint = new Uri("http://jaeger:4318/v1/traces");
                options.Protocol = OpenTelemetry.Exporter.OtlpExportProtocol.HttpProtobuf;
            });
    });

// Register Tracer
builder.Services.AddSingleton<Tracer>(sp => 
    sp.GetRequiredService<TracerProvider>().GetTracer("MyController"));

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
        for (int i = 1; i <= 3; i++)
        {
            command.CommandText = $@"
                INSERT INTO MyTable (Name) VALUES ('Initial Data {i}');
            ";
            command.ExecuteNonQuery();
        }
    }
}
