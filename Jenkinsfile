pipeline {
  agent { label 'aws-minikube' }
  options { disableConcurrentBuilds() }
  environment {
    KUBECTL = 'minikube kubectl --'
    IMAGE   = "web-uoc:${env.BRANCH_NAME}-${env.BUILD_NUMBER}"
  }
  stages {
    stage('Construir imagen') {
      steps { sh 'docker build -t $IMAGE .' }
    }
    stage('Cargar en Minikube') {
      steps { sh 'minikube image load $IMAGE' }
    }
    stage('Desplegar en pruebas') {
      when { branch 'dev' }
      steps {
        sh '''
          set -eu
          $KUBECTL create namespace web-test --dry-run=client -o yaml | $KUBECTL apply -f -
          sed "s|image: web-uoc:.*|image: $IMAGE|" k8s/deployment.yaml | $KUBECTL apply -n web-test -f -
          sed "s/nodePort: 30080/nodePort: 30090/" k8s/service.yaml | $KUBECTL apply -n web-test -f -
          $KUBECTL rollout status deployment/web-uoc -n web-test --timeout=180s
          for i in $(seq 1 20); do curl -sS --connect-timeout 5 --max-time 10 --retry 3 --retry-delay 2 -D - -o /dev/null http://$(minikube ip):30090/ | grep -i '^X-Pod-Name'; done
        '''
      }
    }
    stage('Desplegar en producción') {
      when { branch 'main' }
      steps {
        sh '''
          set -eu
          $KUBECTL create namespace web-prod --dry-run=client -o yaml | $KUBECTL apply -f -
          sed "s|image: web-uoc:.*|image: $IMAGE|" k8s/deployment.yaml | $KUBECTL apply -n web-prod -f -
          $KUBECTL apply -n web-prod -f k8s/service.yaml
          $KUBECTL rollout status deployment/web-uoc -n web-prod --timeout=180s
          for i in $(seq 1 20); do curl -sS --connect-timeout 5 --max-time 10 --retry 3 --retry-delay 2 -D - -o /dev/null http://$(minikube ip):30080/ | grep -i '^X-Pod-Name'; done
        '''
      }
    }
  }
}
