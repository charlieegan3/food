import { CFP_ALLOWED_PATHS } from './constants';
import { getCookieKeyValue } from './utils';
import { getTemplate } from './template';

export async function onRequest(context: {
  request: Request;
  next: () => Promise<Response>;
  env: { CFP_PASSWORD?: string };
}): Promise<Response> {
  const { request, next, env } = context;
  const { pathname, searchParams } = new URL(request.url);
  const { error } = Object.fromEntries(searchParams);
  const cookie = request.headers.get('cookie') || '';
  const cookieKeyValue = await getCookieKeyValue(env.CFP_PASSWORD);

  if (
    cookie.includes(cookieKeyValue) ||
    CFP_ALLOWED_PATHS.includes(pathname) ||
    !env.CFP_PASSWORD
  ) {
    return await next();
  } else {
    // No cookie or incorrect hash in cookie. Redirect to login.
    return new Response(getTemplate({ withError: error === '1', url: request.url }), {
      headers: {
        'content-type': 'text/html'
      }
    });
  }
}
