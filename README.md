# Who's Boosted?

This is a Go server to estimate which (if any) of your Steam friends are boosted. It's a rewrite of an old Flask project, with two major improvements: first, it uses Go and Gin instead of Python and Flask; second, it uses the OpenDota API to get winrates by region, rather than slamming the Steam API.

## Usage

Go [here](https://dota-boost-detective.herokuapp.com/) to try it out in the wild.

For a single account, hit `/api/scan/:steam_id` to get a percentage estimate of whether the account with that steam id is boosted.
