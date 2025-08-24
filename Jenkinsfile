pipeline {
    agent { label 'terraform-agent' }

    environment {
        TF_DIR = 'terraform'
        TF_VAR_postgres_db = 'devops'
        TF_VAR_postgres_user = 'postgres'
        TF_VAR_postgres_host = 'db'
        TF_VAR_postgres_port = 5432
    }

    triggers { pollSCM('H/5 * * * *') }

    stages {

        stage('Checkout') {
            steps {
                git url: "https://github.com/Dzhodddi/devops", branch: 'main'
            }
        }

        stage('Build Go App') {
            steps {
                script {
                    dir("${env.WORKSPACE}") {
                        sh 'go mod tidy'
                        sh 'go build -o ./bin/main ./cmd/api'
                    }
                }
            }
        }

        stage('Test') {
            steps {
                script {
                    dir("${env.WORKSPACE}") {
                        sh '''
                            mkdir -p reports
                                             go test -v ./cmd/api 2>&1 | go run github.com/jstemmer/go-junit-report/v2@latest > reports/report.xml
                        '''
                    }
                }
            }
        }

        stage('Terraform Plan') {
            steps {
                dir("${env.WORKSPACE}/${TF_DIR}") {
                        sh 'terraform init'
                        sh 'terraform plan -out=tfplan'
                    }
                archiveArtifacts artifacts: "${TF_DIR}/tfplan", fingerprint: true
            }
        }

        stage('Terraform Apply') {
            steps {
                dir("${env.WORKSPACE}/${TF_DIR}") {
                    sh 'terraform apply -input=false -auto-approve -target=docker_container.postgres -target=docker_container.app -target=docker_container.frontend'
                }
            }
        }




    }

    post {
        success { echo 'Pipeline completed successfully!' }
        failure { echo 'Pipeline failed!' }
        always { junit 'reports/report.xml' }
    }
}
