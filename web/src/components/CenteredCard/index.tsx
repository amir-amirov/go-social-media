import { Card, CardContent, Container, Typography } from "@mui/material";

interface CenteredCardProps {
  title: string;
  subtitle?: string;
  children: React.ReactNode;
}

const CenteredCard: React.FC<CenteredCardProps> = ({
  title,
  subtitle,
  children,
}) => (
  <Container maxWidth="xs" sx={{ minHeight: "100vh", display: "flex", alignItems: "center" }}>
    <Card
      sx={{
        width: "100%",
        borderRadius: 3,
        boxShadow: 3,
        px: 2,
        py: 4,
      }}
    >
      <CardContent>
        <Typography variant="h4" fontWeight={600} align="center" gutterBottom>
          {title}
        </Typography>
        {subtitle && (
          <Typography variant="body1" color="text.secondary" align="center" mb={2}>
            {subtitle}
          </Typography>
        )}
        {children}
      </CardContent>
    </Card>
  </Container>
);
export default CenteredCard;