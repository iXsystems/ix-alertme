# $FreeBSD$

PORTNAME=	ix-alertme
PORTVERSION=	%%CHGVERSION%%
CATEGORIES=	misc

MAINTAINER=	ken@ixsystems.com
COMMENT=	Plugin-based alert notification system

LICENSE=	BSD2CLAUSE

USE_GITHUB=	yes
GH_ACCOUNT=	ixsystems
GH_PROJECT=	ix-alertme
GH_TAGNAME=	%%GHTAG%%

OPTIONS_DEFINE=	SMTPMAIL
OPTIONS_SUB=		yes
SMTPMAIL_DESC=	Bundle the smtp-email plugin into the package

.include <bsd.port.options.mk>

NO_BUILD=		yes

do-install:
	${INSTALL_PROGRAM} ${WRKSRC}/src-go/ix-alertme/ix-alertme ${STAGEDIR}${PREFIX}/bin/ix-alertme

post-install:
.if ${PORT_OPTIONS:MSMTPMAIL}
	${MKDIR} ${STAGEDIR}/${PREFIX}/${PORTNAME}/plugins/smtp-email
	(cd ${WRKSRC}/provider-plugins/smtp-email && \
	${INSTALL_PROGRAM} smtp-email ${STAGEDIR}/${PREFIX}/${PORTNAME}/plugins/smtp-email/smtp-email)
	(cd ${WRKSRC}/provider-plugins/smtp-email && \
	${CP} manifest.json ${STAGEDIR}/${PREFIX}/${PORTNAME}/plugins/smtp-email/manifest.json)
.endif

.include <bsd.port.mk>
