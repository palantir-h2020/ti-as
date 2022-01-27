#/bin/sh

# Replace server configuration (port)
sed -i "s/server_port: 8100/server_port: ${PORT}/g" "config/config.yaml"

# Replace database configuration (instance, host, port, username, password, db name)
sed -i "s/db_instance: \"redis\"/db_instance: \"${DB_INSTANCE}\"/g" "config/config.yaml"
sed -i "s/db_host: \"localhost\"/db_host: \"${DB_HOST}\"/g" "config/config.yaml"
sed -i "s/db_port: 6379/db_port: ${DB_PORT}/g" "config/config.yaml"
sed -i "s/db_username: \"\"/db_username: \"${DB_USERNAME}\"/g" "config/config.yaml"
sed -i "s/db_password: \"\"/db_password: \"${DB_PASSWORD}\"/g" "config/config.yaml"
sed -i "s/db_name: \"0\"/db_name: \"${DB_NAME}\"/g" "config/config.yaml"

./PalantirIpAnonymizationService "config/config.yaml"