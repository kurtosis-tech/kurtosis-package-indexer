def run(plan):
    plan.add_service(
        name = "indexer",
        config = ServiceConfig(
            image = ImageBuildSpec("./server"),
            ports = {
                "http-api": PortSpec(number = 9770, application_protocol = "http")
            },
            public_ports = {
                "http-api": PortSpec(number = 9770, application_protocol = "http")
            }
        )
    )
