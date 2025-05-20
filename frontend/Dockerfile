FROM ubuntu:22.04 AS build

# Install required packages
RUN apt-get update && apt-get install -y \
    curl \
    git \
    unzip \
    xz-utils \
    zip \
    libglu1-mesa \
    && rm -rf /var/lib/apt/lists/*

# Install Flutter
RUN git clone https://github.com/flutter/flutter.git /flutter
ENV PATH="/flutter/bin:${PATH}"
RUN flutter channel stable
RUN flutter upgrade
RUN flutter doctor

# Set up workspace
WORKDIR /app
COPY . .

# Get dependencies and build web
RUN cd app && flutter pub get
RUN cd app && flutter build web --release

# Production stage
FROM nginx:alpine
COPY --from=build /app/app/build/web /usr/share/nginx/html

# Expose port
EXPOSE 80

CMD ["nginx", "-g", "daemon off;"] 