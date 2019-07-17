# $FreeBSD$

PORTNAME=	ix-alertme
PORTVERSION=	%%CHGVERSION%%
CATEGORIES=	misc

MAINTAINER=	ken@ixsystems.com
COMMENT=	Plugin-based alert notification system

LICENSE=	BSD2CLAUSE

USES=	go:modules

USE_GITHUB=	yes
GH_ACCOUNT=	ixsystems
GH_PROJECT=	ix-alertme
GH_TAGNAME=	%%GHTAG%%

WRKSRC_SUBDIR=	src-go/ix-alertme

OPTIONS_DEFINE=	SMTPMAIL
OPTIONS_SUB=		yes
SMTPMAIL_DESC=	Bundle the smtp-email plugin into the package

.include <bsd.port.options.mk>

post-install:
.if ${PORT_OPTIONS:MSMTPMAIL}
	${MKDIR} ${STAGEDIR}/${PREFIX}/${portname}/plugins/smtp-mail
	(cd ${WRKSRC}/provider-plugins/smtp-email &&
	${COPYTREE_BIN} smtp-email ${STAGEDIR}/${PREFIX}/${portname}/plugins/smtp-mail/smtp-mail)
	(cd ${WRKSRC}/provider-plugins/smtp-email &&
	${COPYTREE_SHARE} manifest.json ${STAGEDIR}/${PREFIX}/${portname}/plugins/smtp-mail/manifest.json)
.endif

.include <bsd.port.mk>
