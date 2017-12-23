---
title: "Riak Core: Quickstart"
date: 2013-01-13
url: "/2013/1/13/riak-core-quickstart/"
---


First Retrieve the Rebar Templates and put them in ~/.rebar/templates

<div><div class="syntaxhighlighter  bash" id="highlighter_831319"><div class="toolbar"><span>[?](#)</span></div><table border="0" cellpadding="0" cellspacing="0"><tbody><tr><td class="gutter"><div class="line number1 index0 alt2">1</div><div class="line number2 index1 alt1">2</div><div class="line number3 index2 alt2">3</div></td><td class="code"><div class="container"><div class="line number1 index0 alt2">`<code class="bash plain"><a href="http://git-scm.com/book">git</a><a href="http://git-scm.com/book/en/Git-Basics-Getting-a-Git-Repository">clone</a><a href="https://github.com/rzezeski/rebar_riak_core">git://github.com/rzezeski/rebar_riak_core.git</a>`</div><div class="line number2 index1 alt1">`<code class="bash functions"><a href="http://unixhelp.ed.ac.uk/CGI/man-cgi?mkdir">mkdir</a>``<code class="bash plain"><a href="http://unixhelp.ed.ac.uk/CGI/man-cgi?mkdir">-p</a> ~/.rebar``<code class="bash plain">/templates`</div><div class="line number3 index2 alt2">`<code class="bash functions"><a href="http://unixhelp.ed.ac.uk/CGI/man-cgi?cp">cp</a>``<code class="bash plain">rebar_riak_core/* ~/.rebar``<code class="bash plain">/templates`</div></div></td></tr></tbody></table></div></div>You’ll also need [Rebar](https://github.com/basho/rebar) and [Erlang](). I used [homebrew](http://mxcl.github.com/homebrew/) to get them:

[bash]
 brew install erlang rebar
 [/bash]

Now that you’re set up, you can create a new multinode app from the templates.
 Replace myapp with whatever you want to call your new app in the following command:

<div><div class="syntaxhighlighter  bash" id="highlighter_654785"><div class="toolbar"><span>[?](#)</span></div><table border="0" cellpadding="0" cellspacing="0"><tbody><tr><td class="gutter"><div class="line number1 index0 alt2">1</div><div class="line number2 index1 alt1">2</div></td><td class="code"><div class="container"><div class="line number1 index0 alt2">`<code class="bash functions">cd``<code class="bash plain">~``<code class="bash plain">/where/I/want/my/app`</div><div class="line number2 index1 alt1">`<code class="bash plain"><a href="https://github.com/basho/rebar">rebar create</a> template=riak_core_multinode appid=myapp nodeid=myapp`</div></div></td></tr></tbody></table></div></div>You can now `make` and run the code.
 [bash]
 make rel
 ./rel/myapp/bin/myapp console
 [/bash]

If all goes well you should be sitting in a console for your app.
 Run
 [erlang]
 myapp:ping().
 [/erlang]
 to ping a random vnode.

##### Links

- [try try try](https://github.com/rzezeski/try-try-try/tree/master/2011/riak-core-first-multinode)
 Post that goes into more detail about the process.
- [Ping]()
 A Future Post detailing the `myapp:ping().` method.