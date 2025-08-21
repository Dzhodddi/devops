pipeline {
  agent { label 'terraform-agent' }

  environment {
    TF_DIR = 'terraform'
    TF_STATE_DIR = '/terraform_state'
  }

  triggers {
    pollSCM('H/1 * * * *') // Poll SCM every 1 minute
  }

  stages {

    stage('Checkout') {
      steps {
        git url: "https://github.com/Dzhodddi/devops", branch: 'main'
      }
    }

    stage('Building...') {
      steps {
        sh 'go mod download'
        sh 'go build -o ./bin/main ./cmd/api'
      }
    }

    stage('Test') {
      steps {
        echo 'Validating terraform...'
        dir("${TF_DIR}") {
          sh 'terraform validate'
        }
        echo 'Testing go app...'
        sh 'go test ./...'
      }
    }

    stage('Deploy') {
      steps {
        echo 'Deploying with Terraform...'
        dir("${TF_DIR}") {
          sh """terraform init -backend-config="path=${TF_STATE_DIR}/terraform.tfstate""""
          sh 'terraform plan -out=tfplan'
          sh 'terraform apply -auto-approve tfplan'
        }
      }
    }

  }

  post {
    success {
      echo 'Pipeline completed successfully!'
    }
    failure {
      echo 'Pipeline failed!'
    }
  }
}
