class ExistsService:
    def __init__(self, client):
        self.client = client



    def exists_post(self, data, headers=None, query_params=None, content_type="application/json"):
        """
        Check existance of a list of keys on the remote storage
        It is method for POST /exists
        """
        uri = self.client.base_url + "/exists"
        return self.client.post(uri, data, headers, query_params, content_type)
