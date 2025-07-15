13/07/2025
# General Principles
🤏🏾 Make changes small, make them often
Small and incremental changes make it easier to spot bugs, are faster to review, and safer to roll back.

🌈 Extensible, not flexible systems
Build extensible systems that excel at specific tasks today and can be expanded later, rather than flexible overly generic systems designed to anticipate multiple future needs. Flexibility often leads to unnecessary complexity by pre-emptively solving hypothetical needs, whereas extensibility defers complexity until we need it. This aligns the design more closely with actual requirements, and minimises guesswork.

One caveat to this is that we would, for example, encourage adding a country code to a data model of a backend service, which might be deemed “flexible”. Doing so however would prevent a time consuming and expensive migration in the future, so do consider when adding a little flexibility comes at very low cost but high potential gain.

🚪 Consider reversibility of decisions
We move quickly at Monzo, and our default should be to make good decisions at high velocity, rather than slow, perfect decisions. However, not all decisions are created equal. When you need to make a decision that is difficult to reverse or has long-lasting impact, take the time to think it through carefully.

⚡ Don’t optimise prematurely
If you can’t show with data it’s a bottleneck, don’t optimise it. Correctness and readability is nearly always more important than performance. 

🧰 Technical debt is a useful tool
Technical debt, when managed thoughtfully, is a strategic tool that allows us to accelerate delivery. Rather than viewing it as a purely negative, we treat it as a deliberate investment that, like financial debt, requires careful consideration and planning for repayment. 

👾 Build systems to be read and debugged by humans
If code is hard to understand, it’s probably too complex. Systems we can’t reason about are especially prone to bugs.
Assume that your reader doesn’t work in your squad and has no context on this system. No-one should need to dive into the implementation of a service to understand its external interface.
