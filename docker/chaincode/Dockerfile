FROM hyperledger/fabric-ccenv

# Install gmp
RUN sudo apt-get update && sudo apt-get install -y libgmp-dev flex bison

# Install pbc
ADD https://crypto.stanford.edu/pbc/files/pbc-0.5.14.tar.gz /root
RUN tar -xvzf /root/pbc-0.5.14.tar.gz \
    && mv /pbc-0.5.14 /root/pbc-0.5.14 \
    && rm /root/pbc-0.5.14.tar.gz \
    && cd /root/pbc-0.5.14 \
    && ./configure \
    && make \
    && sudo make install \
    && sudo ldconfig

# Install pbc go binder
RUN go get github.com/Nik-U/pbc \
    && go install github.com/Nik-U/pbc
