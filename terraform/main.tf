
resource "docker_image" "postgres" {
  name = "postgres:16-alpine"
}

resource "docker_network" "app_network" {
  name = "app_network"
}

resource "docker_network" "jenkins_network" {
  name = "jenkins_network"
}

resource "docker_container" "postgres" {
  name  = "postgres"
  image = docker_image.postgres.name
  restart = "always"
  env = [
    "POSTGRES_USER=${var.postgres_user}",
    "POSTGRES_PASSWORD=${var.postgres_password}",
    "POSTGRES_DB=${var.postgres_db}",
    "POSTGRES_HOST=${var.postgres_host}",
    "POSTGRES_PORT=${var.postgres_port}"
  ]

  ports {
    internal = 5432
    external = 5432
  }

  volumes {
    host_path      = abspath("${path.module}/../pgdata")
    container_path = "/var/lib/postgresql/data"
  }

  networks_advanced {
    name = docker_network.app_network.name
  }
}

resource "docker_image" "app" {
  name = "devops-app"
  build {
    context = "../"
    dockerfile = "Dockerfile"
  }
}

resource "docker_container" "app" {
  image = docker_image.app.name
  name  = "app"
  restart = "always"
  depends_on = [docker_container.postgres]
  env = [
    "email=${var.email}",
    "ENV=development",
    "HTTP_ADDR=0.0.0.0:3000",
    "API_URL=0.0.0.0:3000",
    "POSTGRES_USER=${var.postgres_user}",
    "POSTGRES_PASSWORD=${var.postgres_password}",
    "POSTGRES_DB=${var.postgres_db}",
    "POSTGRES_HOST=${var.postgres_host}",
    "POSTGRES_PORT=${var.postgres_port}",
    "DB_ADDR=postgres://${var.postgres_user}:${var.postgres_password}@${docker_container.postgres.name}:5432/${var.postgres_db}?sslmode=disable"
  ]
  ports {
    internal = 3000
    external = 3000
  }

  command = [
    "./devops"
  ]
  networks_advanced {
    name = docker_network.app_network.name
  }
}

resource "docker_image" "frontend" {
  name = "frontend"
  build {
    context = "../web-src"
    dockerfile = "Dockerfile"
  }
}

resource "docker_container" "frontend" {
  image = docker_image.frontend.name
  name  = "frontend"
  restart = "always"
  depends_on = [docker_container.app]

  ports {
    internal = 5173
    external = 5173
  }

  networks_advanced {
    name = docker_network.app_network.name
  }
}

resource "docker_image" "jenkins" {
  name = "jenkins/jenkins:lts"
}

resource "docker_container" "jenkins_master" {
  name  = "jenkins-master"
  image = docker_image.jenkins.image_id

  ports {
    internal = 4001
    external = 4001
  }

  volumes {
    host_path      = "/var/run/docker.sock"
    container_path = "/var/run/docker.sock"
  }

  env = [
    "JENKINS_OPTS=--httpPort=4001"
  ]

  restart = "always"
  networks_advanced {
    name = docker_network.jenkins_network.name
  }
}
