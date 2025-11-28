<?php

namespace PHPSTORM_META {

    //
    // Relay
    //

    registerArgumentsSet('goridge_relay_socket_type',
        \Spiral\Goridge\SocketRelay::SOCK_TCP,
        \Spiral\Goridge\SocketRelay::SOCK_UNIX,
    );

    expectedArguments(\Spiral\Goridge\SocketRelay::__construct(), 2,
        argumentsSet('goridge_relay_socket_type'));

    //
    // RPC
    //

    registerArgumentsSet('goridge_rpc_options_json',
        \JSON_BIGINT_AS_STRING,
        \JSON_INVALID_UTF8_IGNORE,
        \JSON_INVALID_UTF8_SUBSTITUTE,
        \JSON_OBJECT_AS_ARRAY,
        \JSON_THROW_ON_ERROR,
    );

    registerArgumentsSet('goridge_rpc_options_msgpack',
        \MessagePack\UnpackOptions::BIGINT_AS_DEC,
        \MessagePack\UnpackOptions::BIGINT_AS_GMP,
        \MessagePack\UnpackOptions::BIGINT_AS_STR,
    );

    expectedArguments(\Spiral\Goridge\RPC\RPCInterface::call(), 2,
        argumentsSet('goridge_rpc_options_json'));
    expectedArguments(\Spiral\Goridge\RPC\RPC::call(), 2,
        argumentsSet('goridge_rpc_options_json'));

    expectedArguments(\Spiral\Goridge\RPC\RPCInterface::call(), 2,
        argumentsSet('goridge_rpc_options_msgpack'));
    expectedArguments(\Spiral\Goridge\RPC\RPC::call(), 2,
        argumentsSet('goridge_rpc_options_msgpack'));

    override(\Spiral\Goridge\RPC\RPCInterface::call(), map(['' => '@']));
    override(\Spiral\Goridge\RPC\RPC::call(), map(['' => '@']));

    //
    // RPC Methods
    //

    registerArgumentsSet('goridge_rpc_methods_informer',
        'informer.Workers',
        'informer.List',
    );

    expectedArguments(\Spiral\Goridge\RPC\RPCInterface::call(), 0,
        argumentsSet('goridge_rpc_methods_informer'));
    expectedArguments(\Spiral\Goridge\RPC\RPC::call(), 0,
        argumentsSet('goridge_rpc_methods_informer'));
}
