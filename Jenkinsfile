podTemplate(yaml: '''
apiVersion: v1
kind: Pod
spec:
  nodeSelector:
    k3s.io/hostname: hp-mini
  containers:
  - name: docker
    image: docker:19.03.1
    command:
    - sleep
    args:
    - 99d
    volumeMounts:
      - name: jenkins-docker-cfg
        mountPath: /root/.docker
    env:
      - name: DOCKER_HOST
        value: tcp://localhost:2375
  - name: docker-daemon
    image: docker:19.03.1-dind
    securityContext:
      privileged: true
    env:
      - name: DOCKER_TLS_CERTDIR
        value: ""
  - name: golang
    image: vikaspogu/gocv
    command: ['cat']
    tty: true
  volumes:
  - name: jenkins-docker-cfg
    projected:
      sources:
      - secret:
          name: regcred
          items:
            - key: .dockerconfigjson
              path: config.json
''') {
  node(POD_LABEL) {
    stage ('build app'){
       checkout scm
         container('golang'){
           sh 'cd `pwd` && go build -o main .'
         }
    }
    stage('build image') {
      container('docker') {
          sh 'cd `pwd` && DOCKER_BUILDKIT=1 docker build -t "docker.io/vikaspogu/rpi-node-cm" .'
      }
    }
    stage('push image') {
      container('docker') {
          sh 'DOCKER_BUILDKIT=1 docker push docker.io/vikaspogu/rpi-node-cm'
      }
    }
  }
}
