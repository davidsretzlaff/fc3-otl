FROM mcr.microsoft.com/dotnet/aspnet:8.0 AS base
WORKDIR /app
EXPOSE 80

FROM mcr.microsoft.com/dotnet/sdk:8.0 AS build
ARG BUILD_CONFIGURATION=Release
WORKDIR /src

# Copy all project files first
COPY ["src/Customer.API/Api.csproj", "Customer.API/"]
COPY ["src/Customer.Domain/Domain.csproj", "Customer.Domain/"]
COPY ["src/Customer.Application/Application.csproj", "Customer.Application/"]
COPY ["src/Customer.Infra.Data/Infra.Data.csproj", "Customer.Infra.Data/"]

# Restore dependencies
RUN dotnet restore "Customer.API/Api.csproj"

# Copy the rest of the source code
COPY src/. .

# Build the project
RUN dotnet build "Customer.API/Api.csproj" -c $BUILD_CONFIGURATION -o /app/build

FROM build AS publish
ARG BUILD_CONFIGURATION=Release
RUN dotnet publish "Customer.API/Api.csproj" -c $BUILD_CONFIGURATION -o /app/publish /p:UseAppHost=false

FROM base AS final
WORKDIR /app
COPY --from=publish /app/publish .
ENTRYPOINT ["dotnet", "Api.dll"]
