output "db_id" {
  description = "Postgres id container"
  value = docker_container.postgres.id
}

output "app_id" {
  description = "App id container"
  value = docker_container.app.id
}

output "frontend_id" {
  description = "Frontend id container"
  value = docker_container.frontend.id
}