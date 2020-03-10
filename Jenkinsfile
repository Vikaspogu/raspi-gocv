pipeline {
  agent {
    kubernetes {
      yamlFile "kubernetes/jenkinsPod.yaml"
    }
  }
  options {
      buildDiscarder(logRotator(daysToKeepStr: "30", numToKeepStr: ""))
      disableConcurrentBuilds()
      timeout(time: 1, unit: "HOURS")
  }
  stages {
    stage("build image") {
      steps{
        checkout scm
        container("docker") {
            sh "docker run --rm --privileged multiarch/qemu-user-static --reset -p yes"
            sh "cd `pwd` && DOCKER_CLI_EXPERIMENTAL=enabled DOCKER_BUILDKIT=1 docker build --platform linux/arm64 -t docker.io/vikaspogu/rpi-node-cm:${env.GIT_COMMIT} ."
        }
      }
    }
    stage("push image") {
      steps{
        container("docker") {
          sh "DOCKER_CLI_EXPERIMENTAL=enabled DOCKER_BUILDKIT=1 docker push docker.io/vikaspogu/rpi-node-cm:${env.GIT_COMMIT}"
        }
      }
    }
    stage("deployment"){
      steps{
        container("kubectl"){
          sh "kubectl set image deployment/rpi-node-cm rpi-node-cm=docker.io/vikaspogu/rpi-node-cm:${env.GIT_COMMIT} --record -n raspi-gocv"
        }
      }
    }
  }
  post {
      always {
          cleanWs()
      }
      success {
          slackSend color: 'good', message: '🚀 Build Success'
      }
      failure {
          slackSend color: 'danger', message: '🔥 Build Failure'
      }
  }
}
