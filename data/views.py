from django.contrib.auth.mixins import LoginRequiredMixin
from django.core.exceptions import SuspiciousOperation
from django.http import JsonResponse
# from django.utils.decorators import method_decorator
# from django.views.decorators.csrf import csrf_exempt
from django.views.generic.base import TemplateView
from django.views.generic.list import ListView

from .models import Quote
from .schemas import QuoteBatchSchema, FetchQuotesPayloadSchema


class Quotes(TemplateView):
    template_name = 'quotes/list.html'


# TODO: extract AdminLoginRequiredMixin
class ApiView(LoginRequiredMixin):
    login_url = '/admin/login'
    response_class = JsonResponse


# @method_decorator(csrf_exempt, name='dispatch')
class FetchStockQuotes(ApiView, ListView):
    http_method_names = ['post']
    paginate_by = 100

    def get_queryset(self):
        if not self.request.body:
            raise SuspiciousOperation()

        schema = FetchQuotesPayloadSchema()
        payload, err = schema.loads(self.request.body.decode())
        if err:
            raise SuspiciousOperation(err)

        self.paginate_by = payload.take

        return Quote.objects.filter(stock_id__in=payload.stock_ids, session__date__in=payload.dates)

    def render_to_response(self, context, **kwargs):
        lookup = {
            'items': context['object_list'],
            'pages': context['paginator'].num_pages,
            'per_page': context['paginator'].per_page,
            'current': context['page_obj'].number
        }
        data, _ = QuoteBatchSchema().dump(lookup)
        return self.response_class(data, safe=False)

    def post(self, request, *args, **kwargs):
        return self.get(request, *args, **kwargs)
