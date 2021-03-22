package com.salesforce.tests.fs.cmd;

public interface CmdCommand<T, V> {

    T execute(V args);

}
