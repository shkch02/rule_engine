pipeline {
    agent any 

        environment {
        // 1. Harbor 및 이미지 정보
        HARBOR_URL       = "shkch.duckdns.org"
        HARBOR_PROJECT   = "rule_engine"
        HARBOR_CREDS_ID  = "harbor-creds"
        KUBE_CREDS_ID = "kubeconfig-creds"

        // 2. SSH 터널링/K8s 접속 정보를 환경 변수로 이동 (def 제거)
        K8S_USER = "server4"
        SSH_HOST = "sangsu02.iptime.org"
        K8S_TARGET_IP = "192.168.0.10" 
        K8S_PORT = "6443"

        // 3. 빌드 및 배포 관련 파일명
        DOCKERFILE = "Dockerfile"
        IMAGE_NAME = "rule-engine"
        DEPLOYMENT_YAML = "rule-engine-deployment.yaml"
    }

    stages {
        stage('checkout'){
            steps {
                checkout scm
            }
        }

        stage('Build Docker Image & Push to Harbor') {
            steps {
                script {
                    withCredentials([usernamePassword(credentialsId: HARBOR_CREDS_ID, usernameVariable: 'HARBOR_USER', passwordVariable: 'HARBOR_PASS')]) {
                        sh """
                            docker build -t ${HARBOR_URL}/${HARBOR_PROJECT}/${IMAGE_NAME}:latest .
                            echo $HARBOR_PASS | docker login ${HARBOR_URL} -u $HARBOR_USER --password-stdin
                            docker push ${HARBOR_URL}/${HARBOR_PROJECT}/${IMAGE_NAME}:latest
                        """
                    }
                }
            }
        }
        
    

        stage('Deploy to kubernetes') {
            steps {
                script{
                    def localport = 8888
                    def KUBECONFIG_PATH
                    def tunnelPid
                    def FULL_IMAGE_PATH = "${env.HARBOR_URL}/${env.HARBOR_PROJECT}/${env.IMAGE_NAME}:latest"
                    // ssh 터널 시작
                    sshagent (['k8s-master-ssh-key']){
                        sh "nohup ssh -o StrictHostKeyChecking=no -N -L ${localport}:${env.K8S_TARGET_IP}:${env.K8S_PORT} ${env.K8S_USER}@${env.SSH_HOST} > /dev/null 2>&1 & echo \$! > tunnel.pid"
                        tunnelPid = readFile('tunnel.pid').trim()
                        sleep 10

                    withCredentials([file(credentialsId: KUBE_CREDS_ID, variable : 'KUBE_CONFIG_FILE')]){
                        sh "sed -i 's|server: .*|server: https://127.0.0.1:${localport}|' $KUBE_CONFIG_FILE"
                        sh "sed -i 's|image:.*${env.IMAGE_NAME}:latest|image: ${env.HARBOR_URL}/${env.HARBOR_PROJECT}/${env.IMAGE_NAME}:latest|g' ${DEPLOYMENT_YAML}"
                        echo "Deploying pod with image tag: ${env.IMAGE_TAG}"
                        sh "KUBECONFIG=${KUBE_CONFIG_FILE} kubectl apply -f ${DEPLOYMENT_YAML}"
                        sh "KUBECONFIG=${KUBE_CONFIG_FILE} kubectl rollout status deployment/rule-engine-deployment --timeout=120s"
                        sh "kill ${tunnelPid} || true" 
                        sh "rm -f tunnel.pid || true"                
                        }
                    }
                }
            }
        }
    }

    post {
        always {
           sh "docker logout ${env.HARBOR_URL}"
        }
    }
}