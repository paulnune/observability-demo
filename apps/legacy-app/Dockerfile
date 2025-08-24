FROM python:3.10-slim

WORKDIR /app
COPY app.py .

# Instala dependências necessárias
RUN pip install flask

# Define variáveis de ambiente padrão
ENV LOG_PATH=/var/log/legacy-app/app.log
ENV LOG_INTERVAL=1.0

# Exponha a porta para o servidor HTTP
EXPOSE 8081

# Comando para iniciar o servidor
CMD ["python", "app.py"]