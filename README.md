# grpc-etcd-service-discovery
grpc-etcd-service-discovery


# etcd Setup

## Install etcd on Ubuntu 18.04 / Ubuntu 16.04

* Download the latest release of etcd on Ubuntu 18.04 / Ubuntu 16.04:

          sudo apt -y install wget
          export RELEASE="3.3.13"
          wget https://github.com/etcd-io/etcd/releases/download/v${RELEASE}/etcd-v${RELEASE}-linux-amd64.tar.gz
 
* Extract downloaded archive file.

          tar xvf etcd-v${RELEASE}-linux-amd64.tar.gz
         
* Change to new file directory

          cd etcd-v${RELEASE}-linux-amd64
          
* Move etcd and etcdctl binary files to /usr/local/bin directory.

          sudo mv etcd etcdctl /usr/local/bin
          
* Confirm version

          $ etcd --version
           etcd Version: 3.3.13
           Git SHA: 98d3084
           Go Version: go1.10.8
           Go OS/Arch: linux/amd64
           
           
* Create Etcd configuration file and data directory

          sudo mkdir -p /var/lib/etcd/
          sudo mkdir /etc/etcd
          
* Create etcd system user

          sudo groupadd --system etcd
          sudo useradd -s /sbin/nologin --system -g etcd etcd
          
* Set /var/lib/etcd/ directory ownership to etcd user.
          
          sudo chown -R etcd:etcd /var/lib/etcd/
          
* Configure Systemd and start etcd service

* Create a new systemd service file for etcd.

          sudo vim /etc/systemd/system/etcd.service
          
  Paste below data into the file.
        
          [Unit]
          Description=etcd key-value store
          Documentation=https://github.com/etcd-io/etcd
          After=network.target

          [Service]
          User=etcd
          Type=notify
          Environment=ETCD_DATA_DIR=/var/lib/etcd
          Environment=ETCD_NAME=%m
          ExecStart=/usr/local/bin/etcd
          Restart=always
          RestartSec=10s
          LimitNOFILE=40000

          [Install]
          WantedBy=multi-user.target
          
* Reload systemd service and start etcd on Ubuntu 18,04 / Ubuntu 16,04

          sudo systemctl  daemon-reload
          sudo systemctl  start etcd.service
          
* Check service status:

          $ sudo systemctl  status etcd.service
           ● etcd.service - etcd key-value store
              Loaded: loaded (/etc/systemd/system/etcd.service; disabled; vendor preset: enabled)
              Active: active (running) since Sat 2019-01-05 00:54:20 EAT; 23s ago
                Docs: https://github.com/etcd-io/etcd
            Main PID: 8792 (etcd)
               Tasks: 13 (limit: 4915)
              CGroup: /system.slice/etcd.service
                      └─8792 /usr/local/bin/etcd
           Ama 05 00:54:20 mynix etcd[8792]: 8e9e05c52164694d received MsgVoteResp from 8e9e05c52164694d at term 2
           Ama 05 00:54:20 mynix etcd[8792]: 8e9e05c52164694d became leader at term 2
           Ama 05 00:54:20 mynix etcd[8792]: raft.node: 8e9e05c52164694d elected leader 8e9e05c52164694d at term 2
           Ama 05 00:54:20 mynix etcd[8792]: setting up the initial cluster version to 3.3
           Ama 05 00:54:20 mynix etcd[8792]: set the initial cluster version to 3.3
           Ama 05 00:54:20 mynix etcd[8792]: enabled capabilities for version 3.3
           Ama 05 00:54:20 mynix etcd[8792]: ready to serve client requests
           Ama 05 00:54:20 mynix etcd[8792]: published {Name:5fbf3d068d6c491eb687a7a427fc2263 ClientURLs:[http://localhost:2379]} to cluster cdf818194e3a8c32
           Ama 05 00:54:20 mynix systemd[1]: Started etcd key-value store.
           Ama 05 00:54:20 mynix etcd[8792]: serving insecure client requests on 127.0.0.1:2379, this is strongly discouraged!

* The service will start on localhost address port 2379

         $ ss -tunelp | grep 2379
         tcp   LISTEN 0 128 127.0.0.1:2379 0.0.0.0:*  uid:998 ino:72981 sk:45c <-> 

        $ etcdctl member list
         8e9e05c52164694d: name=5fbf3d068d6c491eb687a7a427fc2263 peerURLs=http://localhost:2380 clientURLs=http://localhost:2379 isLeader=true
