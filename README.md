# gootstrap

Bootstraps Go projects, Birdie style

# Usage

Install using:

```sh
go install github.com/birdie-ai/gootstrap@latest
```

If you checked the project locally:

```sh
go install .
```

Create an empty repository for your project and clone it locally.
Inside the root of the cloned repository run:

```bash
gootstrap -group <service_group> -name <service_name>
```

Then run (we don't provide a go.sum file, so you need to generate one):

```
make mod
```

For more details:

```bash
gootstrap --help
```
