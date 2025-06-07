#!/bin/bash

echo "🧪 Testando configuração OpenTelemetry Collector para logs..."

# Cria logs de teste no formato correto
echo '{"time":"2025-06-07T'$(date +%H:%M:%S)'Z","level":"information","msg":"[TEST] OpenTelemetry Collector test log","correlation_id":"test-'$(date +%Y%m%d%H%M%S)'-12345","service":"test"}' >> logs/otl_test.log

echo '{"time":"2025-06-07T'$(date +%H:%M:%S)'Z","level":"error","msg":"[TEST] Error log for testing","correlation_id":"test-'$(date +%Y%m%d%H%M%S)'-67890","service":"test"}' >> logs/otl_test.log

echo '{"time":"2025-06-07T'$(date +%H:%M:%S)'Z","level":"warning","msg":"[TEST] Warning log for testing","correlation_id":"test-'$(date +%Y%m%d%H%M%S)'-54321","service":"test"}' >> logs/otl_test.log

echo "✅ Logs de teste criados em logs/otl_test.log"
echo ""
echo "📊 Para verificar no Grafana:"
echo "1. Acesse http://localhost:3000"
echo "2. Vá para Explore"
echo "3. Selecione Loki como data source"
echo "4. Use as queries:"
echo "   {service=\"test\"}"
echo "   {level=\"error\"}"
echo "   {correlation_id=~\"test-.*\"}"
echo ""
echo "🔍 Para verificar logs do OpenTelemetry Collector:"
echo "   docker logs otlcollector" 